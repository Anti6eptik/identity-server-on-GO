package service

import (
	"fmt"
	"time"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type TokenService struct {}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (srv TokenService) CreateAcessToken(temp struct {
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

func (srv TokenService) CreateRefreshToken(temp struct {
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

func (srv TokenService) GetTokens(r *http.Request) (*jwt.Token, *jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil, fmt.Errorf("нет заголовка")
	}

	AccessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if AccessTokenString == authHeader {
		return nil, nil, fmt.Errorf("нет Bearer")
	}

	RefreshCookieString, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Print(err)
		return nil, nil, err
	}

	AccessToken, err := srv.ParseToken(AccessTokenString)
	if err != nil {
		return nil, nil, err
	}

	RefreshToken, err := srv.ParseToken(RefreshCookieString.Value)
	if err != nil {
		return nil, nil, err
	}

	return AccessToken, RefreshToken, nil
}

func (srv TokenService) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
