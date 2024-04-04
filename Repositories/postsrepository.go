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
				CREATE (p:Post {content: $content, likes: $likes,comments: $comments ,ImageURL: $imageURL})
			 CREATE (u)-[:POSTED]->(p)`,
			map[string]interface{}{
				"username": username,
				"content":  post.Content,
				"likes":    post.Likes,
				"imageURL": post.ImageURL,
				"comments": post.Comments,
			},
		)
		return nil, err
	})

	return err
}

func GetUserPost(driver neo4j.Driver, username string) ([]data.Post, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (u:User {username: $username})-[:POSTED]->(p:Post)
		RETURN p.content AS content, p.likes AS likes, p.comments AS comments, p.ImageURL AS imageURL
	`
	params := map[string]interface{}{"username": username}

	var posts []data.Post
	result, err := session.Run(query, params)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		record := result.Record()

		var content, imageURL string
		var likes int64
		var commentsSlice []string

		if contentValue, ok := record.Get("content"); ok && contentValue != nil {
			content = contentValue.(string)
		}

		if likesValue, ok := record.Get("likes"); ok && likesValue != nil {
			likes = likesValue.(int64)
		}

		if commentsValue, ok := record.Get("comments"); ok && commentsValue != nil {
			commentsInterfaceSlice := commentsValue.([]interface{})
			for _, comment := range commentsInterfaceSlice {
				if commentStr, ok := comment.(string); ok {
					commentsSlice = append(commentsSlice, commentStr)
				}
			}
		}

		if imageURLValue, ok := record.Get("imageURL"); ok && imageURLValue != nil {
			imageURL = imageURLValue.(string)
		}

		post := data.Post{
			Content:  content,
			Likes:    int(likes),
			Comments: commentsSlice,
			ImageURL: imageURL,
		}
		posts = append(posts, post)
	}
	if err = result.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
