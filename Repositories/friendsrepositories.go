package Repositories

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

func AddFriend(driver neo4j.Driver, usernameSent, usernameRecieved string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (u:User {username: $usernameSent})
						 MATCH (u2:User {username: $usernameRecieved})
						 CREATE (u)-[:FRIEND]->(u2)`,
			map[string]interface{}{
				"usernameSent":     usernameSent,
				"usernameRecieved": usernameRecieved,
			},
		)
		return nil, err
	})
	return err
}

func GetFriendsList(driver neo4j.Driver, username string) ([]string, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
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
