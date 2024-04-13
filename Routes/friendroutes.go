package routes

import (
	service "SocialMedia/Service"
	"SocialMedia/middleware"
	"net/http"
)

func FriendRoutes(mux *http.ServeMux, friendService service.FriendsService) {
	mux.Handle("POST /friends", middleware.AuthMiddleware(http.HandlerFunc(friendService.AddFriend)))
	mux.Handle("DELETE /friends", middleware.AuthMiddleware(http.HandlerFunc(friendService.DeleteFriend)))
	mux.Handle("POST /friends/accept", middleware.AuthMiddleware(http.HandlerFunc(friendService.AcceptFriendRequest)))
}
