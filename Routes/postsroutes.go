package routes

import (
	service "SocialMedia/Service"
	"SocialMedia/middleware"
	"net/http"
)

func PostRoutes(mux *http.ServeMux, postService *service.PostService) {
	mux.Handle("POST /posts/create", middleware.AuthMiddleware(http.HandlerFunc(postService.CreatePost)))
	mux.Handle("GET /posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postService.GetUserPosts)))
	mux.Handle("DELETE /posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postService.DeletePost)))
}
