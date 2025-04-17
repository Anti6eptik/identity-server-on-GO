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
	Service             *service.Service
	HashPasswordService *service.HashPasswordService
}

func NewController(service *service.Service, hashPasswordService *service.HashPasswordService) *Controller {
	return &Controller{
		Service:             service,
		HashPasswordService: hashPasswordService,
	}
}

func (c Controller) InfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Important Info!")
	w.WriteHeader(http.StatusOK)
}

func (c Controller) PostRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var temp struct {
		UserName string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		panic(err)
	}

	HashedPassword, err := c.HashPasswordService.HashPassword(temp.Password)

	if err != nil {
		panic(err)
	}
	temp.Password = HashedPassword

	tokens := c.Service.Registration(temp)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models.TokenModel{Token: tokens.AccessToken}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
}

func (c Controller) PostAuthHandler(w http.ResponseWriter, r *http.Request) {
	var temp struct {
		UserName string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		panic(err)
	}

	tokens, err := c.Service.Auth(temp)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models.TokenModel{Token: tokens.AccessToken}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
}

// СДЕЛАНО!
// Реквест не должен уходить из контроллера декодирование на уровне контроллера
// Secure не работает на нашем сервере
// Access токен передаем через json в теле ответа, модель токен: response
// пофиксиить куки, кеод в банскваде
// Перед возвратом модели access токена надо написать w.Header().Set("Content-Type", "application/json")
// Сменить хост на префикс

// НЕ СДЕЛАНО :(
// Добавить больше ретурнов
// убрать пароль из клеймов это свободная информация
// Шифровка пароля хэшем - хэш теперь сохраняться в БД, а вот авторизацию я не сделал
// Даже с такими куками наш проект не работает(
// Acess токен теперь зраниться в теле, атк как показывал руслан, но я не сделал так, чтобы на сервак брал этот токен из тела, а не из заголовков
// У нас 1 большой сервис, а там и работа с токенами и работа с http все такое, его бы раскидать на TokenService и RegService ну и AuthSerivce
