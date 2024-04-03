package main

import (
	service "SocialMedia/Service"
	"SocialMedia/middleware"
	"fmt"
	"log"
	"net/http"
)

func main() {
	userService := service.NewUserService()
	postService := service.NewPostService()

	mux := http.NewServeMux()

	mux.HandleFunc("/users", userService.CreateUser)
	mux.HandleFunc("/login", userService.LoginUser)

	postsMux := http.NewServeMux()
	postsMux.HandleFunc("/list", handleListPosts)

	protectedCreatePostHandler := middleware.AuthMiddleware(http.HandlerFunc(postService.CreatePost))
	postsMux.Handle("/create", protectedCreatePostHandler)

	mux.Handle("/posts/", http.StripPrefix("/posts", postsMux))
	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleListPosts(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Lista de posts")); err != nil {
		log.Printf("Error escribiendo respuesta: %v", err)
	}
}
