package service

import (
	data "SocialMedia/Data"
	"SocialMedia/Repositories"
	"SocialMedia/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PostService struct {
	driver neo4j.Driver
}

func NewPostService() *PostService {
	driver := db.Driver()
	return &PostService{driver: driver}
}

func (s *PostService) CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost data.Post

	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	username := r.Context().Value("username").(string) // Asumimos que el middleware de autenticación ya ha establecido esto

	if err := Repositories.CreatePost(s.driver, username, newPost); err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Otras funciones del servicio PostService, como recuperar posts, etc.
