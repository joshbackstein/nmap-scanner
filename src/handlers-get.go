package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Serve up the actual index here
}

func GetScansHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["host"]
	log.Println("Getting all scans for host \"" + address + "\"")

	getScansHelper(w, address, 0)
}

func GetNumScansHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["host"]
	numScansString := vars["numScans"]
	log.Println("Getting " + numScansString + " scans for host \"" + address + "\"")

	numScans, err := strconv.Atoi(numScansString)
	if err != nil {
		log.Println("Error parsing number of scans to get")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	scanIdString := vars["scanId"]

	scanId, err := strconv.Atoi(scanIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scan, err := getPreviousScan(db, scanId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			log.Println("Error getting previous scan")
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scan); err != nil {
		log.Println(err.Error())
	}
}
