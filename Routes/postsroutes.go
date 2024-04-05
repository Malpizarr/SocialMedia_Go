package routes

import (
	service "SocialMedia/Service"
	"SocialMedia/middleware"
	"net/http"
)

func PostRoutes(mux *http.ServeMux, postService *service.PostService) {
	mux.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(postService.CreatePost)))
	mux.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postService.GetUserPosts)))
}
