package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Serve up the actual index here
}

func GetScansHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["host"]
	log.Println("Getting scans for host \"" + address + "\"")

	getScansHelper(w, address, 0)
}

func GetNumScansHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["host"]
	numScans, err := strconv.Atoi(vars["numScans"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Getting scans for host \"" + address + "\"")

	getScansHelper(w, address, numScans)
}

func getScansHelper(w http.ResponseWriter, address string, numScans int) {
	scans, err := getScansForHost(db, address, numScans)
	if err != nil {
		log.Println("Error getting scans for host \"" + address + "\"")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scans); err != nil {
		log.Println(err.Error())
	}
}

func GetPreviousScanHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["host"]
	timeString := vars["dateTime"]

	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	dateTime, err := time.Parse(time.RFC3339Nano, timeString)
	if err != nil {
		log.Println("Error parsing datetime")
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scan, err := getPreviousScanForHost(db, address, dateTime)
	if err != nil {
		log.Println("Error getting previous scan for host")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scan); err != nil {
		log.Println(err.Error())
	}
}
