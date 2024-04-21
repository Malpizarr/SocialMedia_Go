package Repositories

import (
	data "SocialMedia/Data"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PostsRepository interface {
	CreatePost(username string, post data.Post) error
	GetUserPost(username string) ([]data.Post, error)
	DeletePost(username, postID string) error
	LikePost(username, postID string) error
	GetLikesFromPost(postId string) ([]string, error)
}

type postsRepository struct {
	driver neo4j.Driver
}

func NewPostsRepository(driver neo4j.Driver) PostsRepository {
	return &postsRepository{driver}
}

func (r *postsRepository) CreatePost(username string, post data.Post) error {
	session := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (u:User {username: $username})
             CREATE (p:Post {id: $id, content: $content, likes: $likes, comments: $comments, ImageURL: $imageURL})
             CREATE (u)-[:POSTED]->(p)`,
			map[string]interface{}{
				"username": username,
				"id":       post.ID,
				"content":  post.Content,
				"likes":    post.Likes,
				"comments": post.Comments,
				"imageURL": post.ImageURL,
			},
		)
		return nil, err
	})

	return err
}

func (r *postsRepository) GetUserPost(username string) ([]data.Post, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (u:User {username: $username})-[:POSTED]->(p:Post)
		RETURN p.id AS ID, p.content AS content, p.likes AS likes, p.comments AS comments, p.ImageURL AS imageURL
	`
	params := map[string]interface{}{"username": username}

	var posts []data.Post
	result, err := session.Run(query, params)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		record := result.Record()

		var id string
		var content, imageURL string
		var likes int64
		var commentsSlice []string

		if IDValue, ok := record.Get("ID"); ok && IDValue != nil {
			id = IDValue.(string)
		}

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
			ID:       id,
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

func (r *postsRepository) DeletePost(username, postID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
            MATCH (u:User {username: $username})-[:POSTED]->(p:Post {id: $postID})
            DETACH DELETE p
        `, map[string]interface{}{
			"username": username,
			"postID":   postID,
		})
		if err != nil {
			log.Printf("Error running Cypher query: %v", err)
			return nil, err
		}

		log.Printf("Delete query result: %v", result)

		return nil, nil
	})
	if err != nil {
		log.Printf("Error in write transaction: %v", err)
	}

	return err
}

func (s *postsRepository) LikePost(username, postID string) error {
	session := s.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			`MATCH (p:Post {id: $postID})
      MATCH (u:User {username: $username})
      MERGE (u)-[r:LIKED]->(p)
      ON CREATE SET r.timestamp = timestamp()
      SET p.likes = p.likes + 1`,
			map[string]interface{}{
				"postID":   postID,
				"username": username,
			})
		return nil, err
	})
	return err
}

func (s *postsRepository) GetLikesFromPost(postId string) ([]string, error) {
	session := s.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			`MATCH (p:Post {id: $postId})<-[:LIKED]-(u:User)
             RETURN u.username AS username`,
			map[string]interface{}{
				"postId": postId,
			},
		)
		if err != nil {
			return nil, err
		}

		var likes []string
		for result.Next() {
			record := result.Record()
			if username, ok := record.Get("username"); ok {
				if username, ok := username.(string); ok {
					likes = append(likes, username)
				}
			}
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return likes, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}
