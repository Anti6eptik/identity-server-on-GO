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

func (c Controller) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	array, code := c.Service.Registration(r)
	if code == "201" {
		w.Write(array)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
