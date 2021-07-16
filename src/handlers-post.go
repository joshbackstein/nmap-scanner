package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func PostScanHandler(w http.ResponseWriter, r *http.Request) {
	// Get host from request
	vars := mux.Vars(r)
	address := vars["host"]

	// Get open ports on host
	ports := getOpenPorts(address)

	// Write to DB
	scan, err := writeToDatabase(db, address, ports)
	if err != nil {
		log.Println("Error writing scan to database")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return scan result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(scan); err != nil {
		log.Println(err.Error())
	}
}
