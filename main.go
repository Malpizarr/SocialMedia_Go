package main

import (
	routes "SocialMedia/Routes"
	service "SocialMedia/Service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env encontrado")
	}

	userService := service.NewUserService()
	postService := service.NewPostService()
	friendService := service.NewFriendsService()
	mux := http.NewServeMux()

	routes.AuthRoutes(mux, userService)
	routes.PostRoutes(mux, postService)
	routes.FriendRoutes(mux, friendService)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "temp/template.html")
	})

	tempFileServer := http.FileServer(http.Dir("temp"))
	mux.Handle("/temp/", http.StripPrefix("/temp/", tempFileServer))

	log.Println("Server starting on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
