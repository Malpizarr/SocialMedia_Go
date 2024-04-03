package Repositories

import (
	data "SocialMedia/Data"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func CreatePost(driver neo4j.Driver, username string, post data.Post) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (u:User {username: $username})
			 CREATE (p:Post {content: $content, likes: $likes})
			 CREATE (u)-[:POSTED]->(p)`,
			map[string]interface{}{
				"username": username,
				"content":  post.Content,
				"likes":    post.Likes,
			},
		)
		return nil, err
	})

	return err
}
