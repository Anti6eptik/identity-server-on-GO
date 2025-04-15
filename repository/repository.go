package repository

import (
	"database/sql"

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
}) int {
	result, err := r.DB.Exec("INSERT INTO users (UserName, Password) VALUES ($1)", new_user.UserName, new_user.Password)
	if err != nil {
		panic(err)
	}
	id, _ := result.LastInsertId()

	return int(id)
}
