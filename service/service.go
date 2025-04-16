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

func (srv Service) GetTokens(r *http.Request) (*jwt.Token, *jwt.Token, error) {
	AccessTokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	RefreshCookieString, err := r.Cookie("refresh_token")
	if err != nil {
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

func (srv Service) ParseToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
}

func (srv Service) AuthMiddleware(next http.Handler) http.Handler{
	return http.HandleFunc(func(w http.ResponseWriter, r *http.Request){
		AccessToken, RefreshToken, err := srv.GetTokens(r)
		if err != nil{
			panic(err)
		}
		if AccessToken.Valid{
			next(w, r)
		} else if RefreshToken.Valid{
			var temp struct {
				UserName string
				Password string
			}

			claims, ok := RefreshToken.Claims.(jwt.MapClaims)
			if !ok{
				panic("МЫ НЕ МОЖЕМ ВЫТАЩИТЬ КЛАЙМЫ")
			}

			UserName, ok := claims["UserName"]
			if !ok{
				panic("МЫ НЕ МОЖЕМ ВЫТАЩИТЬ UserName")
			}

			Password, ok := claims["Password"]
			if !ok{
				panic("МЫ НЕ МОЖЕМ ВЫТАЩИТЬ Password")
			}

			temp.UserName = UserName.(string)
			temp.Password = Password.(string)

			AccessToken, _ := srv.CreateAcessToken(temp)
			RefreshToken, _ := srv.CreateRefreshToken(temp)

			w.Header().Set("Authorization", "Bearer "+AccessToken)
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    RefreshToken,
				HttpOnly: true,
				Secure:   true,
				Path:     "/auth/refresh",
				MaxAge:   60 * 60 * 24 * 7,
			})


		} else{
			
		}
	})
}

func (srv Service) Auth(){
	var temp struct{
		UserName string
		Password string
	}
	if srv.repository.Auth()
}