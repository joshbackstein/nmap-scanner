package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

func main() {
	// Variables
	HTTP_PORT := getEnvOrDefault("HTTP_PORT", "8080")
	MYSQL_USERNAME := getEnvOrDefault("MYSQL_USERNAME", "root")
	MYSQL_PASSWORD := getEnvOrDefault("MYSQL_PASSWORD", "password")
	MYSQL_HOST := getEnvOrDefault("MYSQL_HOST", "localhost")
	MYSQL_PORT := getEnvOrDefault("MYSQL_PORT", "3306")
	MYSQL_DB := getEnvOrDefault("MYSQL_DB", "NmapScanner")

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

func getEnvOrDefault(envVariable string, def string) string {
	var value = os.Getenv("HTTP_PORT")
	if len(value) > 0 {
		return value
	}
	return def
}
