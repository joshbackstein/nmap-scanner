package main

import (
	"strconv"
	"time"
)

type PartialScan struct {
	DateTime time.Time
	Ports    []int
}

type Scan struct {
	Host     string
	DateTime time.Time
	Ports    []int
}

func portsToStringSlice(intPorts []int) []string {
	stringPorts := []string{}
	for _, intPort := range intPorts {
		stringPorts = append(stringPorts, strconv.Itoa(intPort))
	}
	return stringPorts
}

func portsToIntSlice(stringPorts []string) ([]int, error) {
	intPorts := []int{}
	for _, stringPort := range stringPorts {
		intPort, err := strconv.Atoi(stringPort)
		if err != nil {
			return nil, err
		}
		intPorts = append(intPorts, intPort)
	}
	return intPorts, nil
}
