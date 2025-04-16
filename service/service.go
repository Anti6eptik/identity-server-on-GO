package service

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"fmt"

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

func (srv Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AccessToken, RefreshToken, err := srv.GetTokens(r)
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

			newAccessToken, err := srv.CreateAcessToken(temp)
			if err != nil {
				http.Error(w, "Failed to create access token", http.StatusInternalServerError)
				return
			}

			newRefreshToken, err := srv.CreateRefreshToken(temp)
			if err != nil {
				http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Authorization", "Bearer "+newAccessToken)
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    newRefreshToken,
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
				MaxAge:   60 * 60 * 24 * 7,
			})

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/auth", http.StatusFound)
	})
}

func (srv Service) Auth(w http.ResponseWriter, r *http.Request) {
	var temp struct {
		UserName string
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		panic(err)
	}
	if srv.repository.Auth(temp) {
		AccessToken, err := srv.CreateAcessToken(temp)
		if err != nil {
			panic(err)
		}
		RefreshToken, err := srv.CreateRefreshToken(temp)
		if err != nil {
			panic(err)
		}
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
