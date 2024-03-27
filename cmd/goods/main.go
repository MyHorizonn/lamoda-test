package main

import (
	"database/sql"
	"fmt"
	"lamoda-test/internal/handler"
	goods "lamoda-test/internal/storage"
	"lamoda-test/internal/storage/postgres"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	var db goods.Storage
	err := godotenv.Load("./././.env")
	if err != nil {
		log.Fatalln(err)
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	dbOp, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	dbOp.SetMaxOpenConns(10)
	db = &postgres.Postgres{Client: dbOp}
	handler.StartServer(db)
}
