package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func connectDB() {
	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
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


