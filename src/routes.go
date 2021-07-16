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
	Route{"GET", "/", GetRootHandler},
	Route{"GET", "/api/scans/{host}", GetScansHandler},
	Route{"GET", "/api/scans/{host}/{numScans}", GetNumScansHandler},
	Route{"GET", "/api/scans/previous/{host}/{dateTime}", GetPreviousScanHandler},
	Route{"POST", "/api/scans/{host}", PostScanHandler},
}
