package main

import (
	service "SocialMedia/Service"
	"SocialMedia/middleware"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No hay archivo .env")
	}
	userService := service.NewUserService()
	postService := service.NewPostService()

	mux := http.NewServeMux()

	mux.HandleFunc("/users", userService.CreateUser)
	mux.HandleFunc("/login", userService.LoginUser)

	mux.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(postService.CreatePost)))
	mux.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postService.GetUserPosts)))

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
