package service

import (
	"SocialMedia/Repositories"
	"SocialMedia/db"
	"SocialMedia/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	driver neo4j.Driver
}

func NewUserService() *UserService {
	driver := db.Driver()
	return &UserService{driver: driver}
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
		Password string `json:"password" validate:"required,min=8"`
		Email    string `json:"email" validate:"required,email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validateUserData(user); err != nil {
		log.Printf("Invalid user data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := Repositories.CreateUser(s.driver, user.Username, string(hashedPassword), user.Email); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	createdUser, err := Repositories.GetUser(s.driver, user.Username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdUser)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *UserService) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := validateLoginCredentials(credentials); err != nil {
		log.Printf("Invalid login credentials: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := Repositories.GetUser(s.driver, credentials.Username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(credentials.Username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return token as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func validateUserData(user struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
},
) error {
	// Validate username
	if len(user.Username) < 4 || len(user.Username) > 20 || !isAlphanumeric(user.Username) {
		return errors.New("Invalid username")
	}

	// Validate password
	if len(user.Password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	// Validate email
	if !isValidEmail(user.Email) {
		return errors.New("Invalid email address")
	}

	return nil
}

func validateLoginCredentials(credentials struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8"`
},
) error {
	// Validate username
	if len(credentials.Username) < 4 || len(credentials.Username) > 20 || !isAlphanumeric(credentials.Username) {
		return errors.New("Invalid username")
	}

	// Validate password
	if len(credentials.Password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	return nil
}

func isAlphanumeric(str string) bool {
	pattern := "^[a-zA-Z0-9]*$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
