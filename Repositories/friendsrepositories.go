package Repositories

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type FriendsRepository interface {
	AddFriend(usernameSent, usernameRecieved string) error
	GetFriendsList(username string) ([]string, error)
	DeleteFriend(usernamesent, usernamereceived string) error
	AcceptFriendRequest(usernameSent, usernameRecieved string) error
}

type friendsRepository struct {
	driver neo4j.Driver
}

func NewFriendsRepository(driver neo4j.Driver) FriendsRepository {
	return &friendsRepository{driver: driver}
}

func (graph *friendsRepository) AddFriend(usernameSent, usernameRecieved string) error {
	session := graph.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (u:User {username: $usernameSent})
			 MATCH (u2:User {username: $usernameRecieved})
			 MERGE (u)-[r:FRIEND]->(u2)
       ON CREATE SET r.acepted = false`,
			map[string]interface{}{
				"usernameSent":     usernameSent,
				"usernameRecieved": usernameRecieved,
			},
		)
		return nil, err
	})
	return err
}

func (graph *friendsRepository) AcceptFriendRequest(usernameSent, usernameRecieved string) error {
	session := graph.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (u:User {username: $usernameSent})-[r:FRIEND]-(u2:User {username: $usernameRecieved})
     SET r.acepted = true`,
			map[string]interface{}{
				"usernameSent":     usernameSent,
				"usernameRecieved": usernameRecieved,
			})
		return nil, err
	})
	return err
}

func (graph *friendsRepository) GetFriendsList(username string) ([]string, error) {
	session := graph.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
        MATCH (u:User {username: $username})-[:FRIEND]-(friend:User)
        RETURN friend.username AS friendUsername
    `

	result, err := session.Run(query, map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, err
	}

	var friends []string
	for result.Next() {
		record := result.Record()
		friendUsername, _ := record.Get("friendUsername")
		friends = append(friends, friendUsername.(string))
	}

	return friends, nil
}

func (graph *friendsRepository) DeleteFriend(usernameSent, usernameRecieved string) error {
	session := graph.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (u:User {username: $usernameSent})-[f:FRIEND]-(u2:User {username: $usernameRecieved})
						 DELETE f`,
			map[string]interface{}{
				"usernameSent":     usernameSent,
				"usernameRecieved": usernameRecieved,
			},
		)
		return nil, err
	})
	return err
}
