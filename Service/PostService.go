package service

import (
	data "SocialMedia/Data"
	"SocialMedia/Repositories"
	"SocialMedia/db"
	"SocialMedia/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PostService struct {
	friendRepo Repositories.FriendsRepository
	driver     neo4j.Driver
}

func NewPostService() *PostService {
	driver := db.Driver()
	return &PostService{driver: driver}
}

func (s *PostService) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error al parsear el formulario multipart: %v", err)
		http.Error(w, "Error al procesar la carga del archivo", http.StatusBadRequest)
		return
	}

	var newPost data.Post
	newPost.ID = uuid.New().String()
	newPost.Content = r.FormValue("content")
	newPost.Likes = 0

	comments := r.FormValue("comments")
	if comments != "" {
		newPost.Comments = strings.Split(comments, ",")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error al leer el archivo: %v", err)
		http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
		return
	}

	containerName := "posts"
	blobURL, err := utils.UploadFileToBlobStorage(containerName, header.Filename, fileBytes)
	if err != nil {
		log.Printf("Error al subir el archivo a Blob Storage: %v", err)
		http.Error(w, "Error al subir el archivo", http.StatusInternalServerError)
		return
	}
	newPost.ImageURL = blobURL

	username := r.Context().Value("username").(string)

	if err := Repositories.CreatePost(s.driver, username, newPost); err != nil {
		log.Printf("Error creando post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Post creado con éxito, imagen almacenada en: " + blobURL)); err != nil {
		log.Printf("Error escribiendo respuesta: %v", err)
	}
}

func (s *PostService) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("id")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	posts, err := Repositories.GetUserPost(s.driver, username)
	if err != nil {
		log.Printf("Error obteniendo posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *PostService) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)

	if err := Repositories.DeletePost(s.driver, username, postID); err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Post deleted successfully")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *PostService) GetFriendsPosts(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	friends, err := s.friendRepo.GetFriendsList(username)
	if err != nil {
		log.Printf("Error obteniendo lista de amigos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var friendsPosts []data.Post
	for _, friend := range friends {
		posts, err := Repositories.GetUserPost(s.driver, friend)
		if err != nil {
			log.Printf("Error obteniendo posts del amigo %s: %v", friend, err)
			continue
		}
		friendsPosts = append(friendsPosts, posts...)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friendsPosts); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
