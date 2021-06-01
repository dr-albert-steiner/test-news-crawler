package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func connectDB() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db = database
}

func disconnectDB(){
	err := db.Close()
	if err != nil {
		log.Println(err.Error())
	}
}


