package repository

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite", "db/database.db")
	if err != nil {
		panic(err)
	}
	return db
}

func (r Repository) Registration(new_user struct {
	UserName string
	Password string
}) error {
	_, err := r.DB.Exec("INSERT INTO users (UserName, Password) VALUES ($1, $2)", new_user.UserName, new_user.Password)
	if err != nil {
		return fmt.Errorf("ошибка базы данных")
	}
	return nil
}

func (r Repository) GetPasswordHash(user struct {
	UserName string
	Password string
}) (string, error){
	var temp struct {
		UserName string
		Password string
	}
	row := r.DB.QueryRow("SELECT * FROM users WHERE UserName=$1", user.UserName)
	var id int
	err := row.Scan(&id, &temp.UserName, &temp.Password)
	if err != nil {
		return "", err
	}
	return temp.Password, nil
}
