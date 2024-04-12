package routes

import (
	service "SocialMedia/Service"
	"SocialMedia/middleware"
	"net/http"
)

func FriendRoutes(mux *http.ServeMux, friendService *service.FriendsService) {
	mux.Handle("POST /friends/add", middleware.AuthMiddleware(http.HandlerFunc(friendService.AddFriend)))
	mux.Handle("DELETE /friends/delete", middleware.AuthMiddleware(http.HandlerFunc(friendService.DeleteFriend)))
}
