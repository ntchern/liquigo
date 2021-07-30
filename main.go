package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/ntchern/liquigo/liquigo"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expected change log file name")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DATABASE")
	schema := os.Getenv("DB_SCHEMA")
	sslMode := os.Getenv("DB_SSLMODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s search_path=%s sslmode=%s",
		host, port, user, password, dbname, schema, sslMode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	fileName := os.Args[1]
	err = liquigo.Update(db, fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
}
