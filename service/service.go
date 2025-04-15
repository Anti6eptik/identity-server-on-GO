package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Sinekod/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (srv Service) Get_json_id(id int) []byte { //любой вывод json id
	var choto = map[string]int{"id": id}
	data, err := json.Marshal(choto)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (srv Service) Registration(r *http.Request) ([]byte, string) {
	var temp struct {
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		return nil, "400"
	}

	id := srv.repository.Registration(temp)

	return srv.Get_json_id(id), "201"
}
