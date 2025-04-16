package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"go.uber.org/dig"

	"Sinekod/controller"
	"Sinekod/repository"
	"Sinekod/service"
)

func main() {
	container := dig.New()

	_ = container.Provide(controller.NewController)
	_ = container.Provide(service.NewService)
	_ = container.Provide(repository.NewRepository)
	_ = container.Provide(repository.NewDB)

	container.Invoke(func(controller *controller.Controller) {
		router := mux.NewRouter()

		router.HandleFunc("/registration", controller.PostRegistrationHandler).Methods("POST")

		router.HandleFunc("/auth", controller.PostAuthHandler).Methods("POST")

		ImportantInfo := router.Host("localhost:8080").Subrouter()
		ImportantInfo.Use(controller.Service.AuthMiddleware)
		ImportantInfo.HandleFunc("/info", controller.InfoHandler).Methods("GET")

		fmt.Println("Server listening...")
		http.ListenAndServe(":8080", router)
	})
}
