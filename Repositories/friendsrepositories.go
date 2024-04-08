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
