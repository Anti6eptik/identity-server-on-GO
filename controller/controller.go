package controller

import (
	"Sinekod/models"
	"Sinekod/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

func (c Controller) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
	w.WriteHeader(http.StatusOK)
}

func (c Controller) PostRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var temp struct {
		UserName string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	HashedPassword, err := c.Service.HashPasswordService.HashPassword(temp.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	temp.Password = HashedPassword

	tokens, err := c.Service.Registration(temp)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		// Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   60 * 60 * 24 * 7,
	})
	json.NewEncoder(w).Encode(models.TokenModel{Token: tokens.AccessToken})
}

func (c Controller) PostAuthHandler(w http.ResponseWriter, r *http.Request) {
	var temp struct {
		UserName string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := c.Service.Auth(temp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		// Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   60 * 60 * 24 * 7,
	})
	json.NewEncoder(w).Encode(models.TokenModel{Token: tokens.AccessToken})
}