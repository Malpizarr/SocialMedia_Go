package service

import (
	"SocialMedia/Repositories"
	"SocialMedia/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

type userService struct {
	userRepo Repositories.UserRepository
}

func NewUserService(userRepo Repositories.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) Register(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
		Password string `json:"password" validate:"required,min=8"`
		Email    string `json:"email" validate:"required,email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "request body invalid", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validateUserData(user); err != nil {
		log.Printf("Datos de usuario invalidos: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hasheando la contrasena: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := s.userRepo.CreateUser(user.Username, string(hashedPassword), user.Email); err != nil {
		if err.Error() == "el username ya está en uso" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			log.Printf("Error creando el usuario: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	createdUser, err := s.userRepo.GetUser(user.Username)
	if err != nil {
		log.Printf("Error al obtener el usuario: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *userService) LoginUser(w http.ResponseWriter, r *http.Request) {
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

	if err := validateLoginCredentials(credentials); err != nil {
		log.Printf("Credenciales de inicio invalidas: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userRepo.GetUser(credentials.Username)
	if err != nil {
		log.Printf("Error al obtener el usuario: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Usuario o contrasena invalidos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(credentials.Password)); err != nil {
		http.Error(w, "Usuario o contrasena invalidos", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(credentials.Username)
	if err != nil {
		log.Printf("Error generando el token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    token,
		"username": credentials.Username,
	}); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func validateUserData(user struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
},
) error {
	if len(user.Username) < 4 || len(user.Username) > 20 || !isAlphanumeric(user.Username) {
		return errors.New("username invalido")
	}

	if len(user.Password) < 8 {
		return errors.New("la contrasena debe tener al menos 8 caracteres")
	}

	if !isValidEmail(user.Email) {
		return errors.New("email invalido")
	}

	return nil
}

func validateLoginCredentials(credentials struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8"`
},
) error {
	if len(credentials.Username) < 4 || len(credentials.Username) > 20 || !isAlphanumeric(credentials.Username) {
		return errors.New("username invalido")
	}

	if len(credentials.Password) < 8 {
		return errors.New("la contrasena debe tener al menos 8 caracteres")
	}

	return nil
}

func isAlphanumeric(str string) bool {
	pattern := "^[a-zA-Z0-9]*$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9_+&*-]+(?:\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,7}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
