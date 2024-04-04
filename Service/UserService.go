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

	if err := Repositories.CreateUser(s.driver, user.Username, string(hashedPassword), user.Email); err != nil {
		log.Printf("Error creando el usuario: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	createdUser, err := Repositories.GetUser(s.driver, user.Username)
	if err != nil {
		log.Printf("Error al obtener el usuario: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdUser)
	if err != nil {
		log.Printf("Error serializando el usuario: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Error escribiendo la respuesta: %v", err)
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

	if err := validateLoginCredentials(credentials); err != nil {
		log.Printf("Credenciales de inicio invalidas: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := Repositories.GetUser(s.driver, credentials.Username)
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
	json.NewEncoder(w).Encode(map[string]string{"token": token})
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
		return errors.New("Email invalido")
	}

	return nil
}

func validateLoginCredentials(credentials struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8"`
},
) error {
	if len(credentials.Username) < 4 || len(credentials.Username) > 20 || !isAlphanumeric(credentials.Username) {
		return errors.New("Ususername invalido")
	}

	if len(credentials.Password) < 8 {
		return errors.New("La contrasena debe tener al menos 8 caracteres")
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
