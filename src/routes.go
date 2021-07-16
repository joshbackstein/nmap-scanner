package main

import (
	"net/http"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{"GET", "/api/hosts/{host}/scans", GetScansHandler},
	Route{"GET", "/api/hosts/{host}/scans/{numScans}", GetNumScansHandler},
	Route{"GET", "/api/scans/{scanId}/previousByHost", GetPreviousScanHandler},
	Route{"POST", "/api/hosts/{host}/scans", PostScanHandler},
}
