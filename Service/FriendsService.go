package service

import (
	"SocialMedia/Repositories"
	"encoding/json"
	"log"
	"net/http"
)

type FriendsService interface {
	AddFriend(w http.ResponseWriter, r *http.Request)
	DeleteFriend(w http.ResponseWriter, r *http.Request)
	AcceptFriendRequest(w http.ResponseWriter, r *http.Request)
}

type friendsService struct {
	FriendRepo Repositories.FriendsRepository
}

func NewFriendsService(fr Repositories.FriendsRepository) FriendsService {
	return &friendsService{fr}
}

func (s *friendsService) AddFriend(w http.ResponseWriter, r *http.Request) {
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
	if err := s.FriendRepo.AddFriend(friendRequest.UsernameSent, friendRequest.UsernameReceived); err != nil {

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

func (s *friendsService) DeleteFriend(w http.ResponseWriter, r *http.Request) {
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
	if err := s.FriendRepo.DeleteFriend(friendRequest.UsernameSent, friendRequest.UsernameReceived); err != nil {
		log.Printf("Error deleting friend: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"message": "Friend deleted"}`)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *friendsService) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
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
	if err := s.FriendRepo.AcceptFriendRequest(friendRequest.UsernameSent, friendRequest.UsernameReceived); err != nil {
		log.Printf("Error accepting friend: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"message": "Friend accepted"}`)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
