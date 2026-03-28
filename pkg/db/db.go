package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := "host=localhost port=5433 user=postgres password=postgres dbname=jobs sslmode=disabled"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	fmt.Println("Connected to Database")
}
