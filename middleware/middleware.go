package middleware

import (
	"net/http"
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		next(w, r)
	}
}