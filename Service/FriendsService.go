package service

import (
	"SocialMedia/Repositories"
	"SocialMedia/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FriendsService struct {
	driver neo4j.Driver
}

func NewFriendsService() *FriendsService {
	driver := db.Driver()
	return &FriendsService{driver: driver}
}

func (s *FriendsService) AddFriend(w http.ResponseWriter, r *http.Request) {
	var friendRequest struct {
		UsernameSent     string `json:"usernamesent"`
		UsernameReceived string `json:"usernamereceived"`
	}
	if err := json.NewDecoder(r.Body).Decode(&friendRequest); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "request body invalid", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := Repositories.AddFriend(s.driver, friendRequest.UsernameSent, friendRequest.UsernameReceived); err != nil {
		log.Printf("Error adding friend: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte(`{"message": "Friend added"}`)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
