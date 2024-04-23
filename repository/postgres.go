package repository

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func InitDB() (*sql.DB, error) {
	godotenv.Load(".env")
	databaseSource := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func New(db *sql.DB) *Postgres {

	return &Postgres{Db: db}
}
