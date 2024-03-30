package main

import (
	"database/sql"
	"fmt"
	"lamoda-test/internal/handler"
	goods "lamoda-test/internal/storage"
	"lamoda-test/internal/storage/postgres"
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	var db goods.Storage
	err := godotenv.Load("./././.env")
	if err != nil {
		log.Fatalln(err)
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), "host.docker.internal", os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	dbOp, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbOp.Close()
	dbOp.SetMaxOpenConns(runtime.NumCPU())
	db = &postgres.Postgres{Client: dbOp}
	handler.StartServer(db)
}
