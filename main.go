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
		log.Print("No hay archivo .env")
	}
	userService := service.NewUserService()
	postService := service.NewPostService()

	mux := http.NewServeMux()

	routes.AuthRoutes(mux, userService)
	routes.PostRoutes(mux, postService)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
