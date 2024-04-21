package service

import (
	data "SocialMedia/Data"
	"SocialMedia/Repositories"
	"SocialMedia/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type PostService interface {
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetUserPosts(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	GetFriendsPosts(w http.ResponseWriter, r *http.Request)
	LikePost(w http.ResponseWriter, r *http.Request)
	GetLikesFromPost(w http.ResponseWriter, r *http.Request)
}

type postService struct {
	friendRepo Repositories.FriendsRepository
	postRepo   Repositories.PostsRepository
}

func NewPostService(pr Repositories.PostsRepository, fr Repositories.FriendsRepository) PostService {
	return &postService{postRepo: pr, friendRepo: fr}
}

func (s *postService) CreatePost(w http.ResponseWriter, r *http.Request) {
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

	if err := s.postRepo.CreatePost(username, newPost); err != nil {
		log.Printf("Error creando post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Post creado con éxito, imagen almacenada en: " + blobURL)); err != nil {
		log.Printf("Error escribiendo respuesta: %v", err)
	}
}

func (s *postService) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("id")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	posts, err := s.postRepo.GetUserPost(username)
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

func (s *postService) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)

	if err := s.postRepo.DeletePost(username, postID); err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Post deleted successfully")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *postService) GetFriendsPosts(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	friends, err := s.friendRepo.GetFriendsList(username)
	if err != nil {
		log.Printf("Error obteniendo lista de amigos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var friendsPosts []data.Post
	for _, friend := range friends {
		posts, err := s.postRepo.GetUserPost(friend)
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

func (s *postService) LikePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostID string `json:"postID"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decodificando el cuerpo de la solicitud: %v", err)
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	if err := s.postRepo.LikePost(username, req.PostID); err != nil {
		log.Printf("Error dando like al post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Post liked successfully")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *postService) GetLikesFromPost(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("postID")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	likes, err := s.postRepo.GetLikesFromPost(postID)
	if err != nil {
		log.Printf("Error obteniendo likes del post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(likes); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
