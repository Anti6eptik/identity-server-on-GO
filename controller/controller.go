package controller

import (
	"Sinekod/service"
	"fmt"
	"net/http"
)

type Controller struct {
	Service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		Service: service,
	}
}

func (c Controller) InfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Important Info!")
	w.WriteHeader(http.StatusOK)
}

func (c Controller) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	AccessToken, RefreshToken, err := c.Service.Registration(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Authorization", "Bearer "+AccessToken)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    RefreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/auth/refresh",
			MaxAge:   60 * 60 * 24 * 7,
		})
	}
}

func (c Controller) AuthHandler(w http.ResponseWriter, r *http.Request) {

}
