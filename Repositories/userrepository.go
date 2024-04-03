package Repositories

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func CreateUser(driver neo4j.Driver, username, password, email string) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"CREATE (u:User {username: $username, password: $password, email: $email}) RETURN u",
			map[string]interface{}{"username": username, "password": password, "email": email},
		)
		return nil, err
	})
	return err
}

func GetUser(driver neo4j.Driver, username string) (map[string]interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	result, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (u:User {username: $username}) RETURN u",
			map[string]interface{}{"username": username},
		)
		if err != nil {
			return nil, err
		}
		if result.Next() {
			record := result.Record()
			if node, ok := record.Values[0].(neo4j.Node); ok {
				return node.Props, nil
			}
		}
		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(map[string]interface{}), nil
}
