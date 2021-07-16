package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	// Variables
	HTTP_PORT := "8080"
	MYSQL_USERNAME := "root"
	MYSQL_PASSWORD := "password"
	MYSQL_HOST := "localhost"
	MYSQL_PORT := "3306"
	MYSQL_DB := "NmapScanner"

	// DB setup
	connectionString := MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/" + MYSQL_DB
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()
	log.Println("Connected to database")

	// HTTP server setup
	router := NewRouter()
	log.Println("Listening for HTTP requests on port " + HTTP_PORT)
	log.Fatal(http.ListenAndServe(":"+HTTP_PORT, router))
}
