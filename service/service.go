package service

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"Sinekod/repository"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("TIMOFEY_NE_LUBIT_GRECHKU")

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (srv Service) Registration(r *http.Request) (string, string, error) {
	var temp struct {
		UserName string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		return "", "", err
	}

	srv.repository.Registration(temp)

	AccessToken, _ := srv.CreateAcessToken(temp)
	RefreshToken, _ := srv.CreateRefreshToken(temp)

	return AccessToken, RefreshToken, nil
}

func (srv Service) CreateAcessToken(temp struct {
	UserName string
	Password string
}) (string, error) {
	AccessCclaims := jwt.MapClaims{
		"UserName": temp.UserName,
		"Password": temp.Password,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessCclaims)

	return AccessToken.SignedString(secretKey)
}

func (srv Service) CreateRefreshToken(temp struct {
	UserName string
	Password string
}) (string, error) {
	RefreshClaims := jwt.MapClaims{
		"UserName": temp.UserName,
		"Password": temp.Password,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	}

	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims)

	return RefreshToken.SignedString(secretKey)
}

func (srv Service) GetTokens(r *http.Request) (string, string, error) {
	accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshCookie.Value, nil
}


