package routes

import (
	service "SocialMedia/Service"
	"net/http"
)

func AuthRoutes(mux *http.ServeMux, userService service.UserService) {
	mux.HandleFunc("/register", userService.Register)
	mux.HandleFunc("/login", userService.LoginUser)
}
