package service

import (
	"Sinekod/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"Sinekod/models"
)

var secretKey = []byte("TIMOFEY_NE_LUBIT_GRECHKU")

type Service struct {
	repository          *repository.Repository
	HashPasswordService *HashPasswordService
	TokenService *TokenService
}

func NewService(repo *repository.Repository, hashPasswordService *HashPasswordService) *Service {
	return &Service{
		repository:          repo,
		HashPasswordService: hashPasswordService,
	}
}

func (srv Service) Registration(temp struct {
	UserName string
	Password string
}) (models.Tokens, error) {
	var tokens models.Tokens
	err := srv.repository.Registration(temp)
	if err != nil {
		return tokens, err
	}

	AccessToken, err := srv.TokenService.CreateAcessToken(temp)
	if err != nil {
		return tokens, err
	}
	RefreshToken, err := srv.TokenService.CreateRefreshToken(temp)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = AccessToken
	tokens.RefreshToken = RefreshToken

	return tokens, nil
}


func (srv Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AccessToken, RefreshToken, err := srv.TokenService.GetTokens(r)
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusFound)
			return
		}

		if AccessToken.Valid {
			next.ServeHTTP(w, r)
			return
		}

		if RefreshToken.Valid {
			claims, ok := RefreshToken.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userName, ok := claims["UserName"].(string)
			if !ok {
				http.Error(w, "Invalid UserName in token", http.StatusUnauthorized)
				return
			}

			password, ok := claims["Password"].(string)
			if !ok {
				http.Error(w, "Invalid Password in token", http.StatusUnauthorized)
				return
			}

			temp := struct {
				UserName string
				Password string
			}{
				UserName: userName,
				Password: password,
			}

			newAccessToken, err := srv.TokenService.CreateAcessToken(temp)
			if err != nil {
				http.Error(w, "Failed to create access token", http.StatusInternalServerError)
				return
			}

			newRefreshToken, err := srv.TokenService.CreateRefreshToken(temp)
			if err != nil {
				http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(models.TokenModel{Token: newAccessToken}); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    newRefreshToken,
				HttpOnly: true,
				//Secure:   true,
				Path:     "/",
				SameSite: http.SameSiteNoneMode,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				MaxAge:   60 * 60 * 24 * 7,
			})

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/auth", http.StatusFound)
	})
}

func (srv Service) Auth(temp struct {
	UserName string
	Password string
}) (models.Tokens, error) {
	var tokens models.Tokens
	RightPassword, err := srv.repository.GetPasswordHash(temp)
	if err != nil {
		return tokens, err
	}
	if srv.HashPasswordService.CheckPasswordHash(temp.Password, RightPassword) {
		var user struct {
			UserName string
			Password string
		}
		user.UserName = temp.UserName
		user.Password = RightPassword
		AccessToken, err := srv.TokenService.CreateAcessToken(user)
		if err != nil {
			return tokens, err
		}
		RefreshToken, err := srv.TokenService.CreateRefreshToken(user)
		if err != nil {
			return tokens, err
		}
		tokens.AccessToken = AccessToken
		tokens.RefreshToken = RefreshToken

		return tokens, nil
	}
	return tokens, fmt.Errorf("ошибка авторизации")
}
