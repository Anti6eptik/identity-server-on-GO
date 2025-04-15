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
		Registration := mux.NewRouter()

		Registration.HandleFunc("/registration", controller.RegistrationHandler).Methods("POST")

		ImportantInfo := mux.NewRouter()

		ImportantInfo.HandleFunc("/info", controller.InfoHandler)

		fmt.Println("Server listening...")
		http.ListenAndServe(":8080", Registration)
		http.ListenAndServe(":8080", ImportantInfo)
	})
}
