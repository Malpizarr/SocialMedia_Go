package main

import (
	"SocialMedia/Repositories"
	routes "SocialMedia/Routes"
	service "SocialMedia/Service"
	"SocialMedia/db"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env encontrado")
	}

	friendrepo := Repositories.NewFriendsRepository(db.Driver())
	postrepo := Repositories.NewPostsRepository(db.Driver())
	userrepo := Repositories.NewUserRepository(db.Driver())

	userService := service.NewUserService(userrepo)
	postService := service.NewPostService(postrepo, friendrepo)
	friendService := service.NewFriendsService(friendrepo)
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
