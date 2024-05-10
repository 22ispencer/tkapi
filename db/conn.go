package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

type Env struct {
	DB *sql.DB
}

func GetEnv() (*Env, error) {
	godotenv.Load()

	server := os.Getenv("DB_SERVER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s", server, user, password, database, port)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return &Env{}, err
	}

	return &Env{DB: db}, nil
}

func (e *Env) Close() {
	e.DB.Close()
}
