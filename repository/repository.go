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
}) {
	_, err := r.DB.Exec("INSERT INTO users (UserName, Password) VALUES ($1, $2)", new_user.UserName, new_user.Password)
	if err != nil {
		panic(err)
	}
}

func (r Repository) Auth(user struct {
	UserName string
	Password string
}) bool {
	var temp struct {
		UserName string
		Password string
	}
	row := r.DB.QueryRow("SELECT * FROM users WHERE UserName=$1, Password=$2", user.UserName, user.Password)
	err := row.Scan(&temp.UserName, &temp.Password)
	if err != nil {
		panic(err)
	}
	if temp == user {
		return true
	} else {
		return false
	}
}
