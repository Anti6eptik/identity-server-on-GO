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


func (c Controller) PostRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	AccessToken, RefreshToken, err := c.Service.Registration(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Authorization", "Bearer "+AccessToken)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    RefreshToken,
			HttpOnly: true,
			// Secure:   false, // Для локального тестирования
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 7,
		})
	}
}


func (c Controller) PostAuthHandler(w http.ResponseWriter, r *http.Request) {
	c.Service.Auth(w, r)
	w.WriteHeader(http.StatusOK)
}


// Реквест не должен уходить из контроллера декодирование на уровне контроллера
// Secure не работает на нашем 
// Access токен передаем через json в теле ответа, модель токен: response
// посмотри код в бан скваде
// EXPIRES - время жизни
// MaxAge - максимальное время жизни
// пофиксиить куки, кеод в банскваде
// Перед возвратом модели access токена надо написать w.Header().Set("Content-Type", "application/json")
// Добавить больше ретурнов
// убрать пароль из клеймов это свободная информация
// Сменить хост на префикс