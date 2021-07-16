package main

import (
	"errors"
	"log"
	"net"
	"os/exec"
	"strings"
	"unicode"
)

var InvalidHostError error = errors.New("Invalid host")

func getOpenPorts(address string) ([]string, error) {
	output, err := runNmap(address)
	if err != nil {
		return nil, err
	}
	return parsePorts(*output), nil
}

func runNmap(address string) (*string, error) {
	if !isValidHost(address) {
		return nil, InvalidHostError
	}

	log.Println("Running Nmap for \"" + address + "\"")
	outputBytes, err := exec.Command("nmap", "--open", "-p0-1000", address).Output()
	if err != nil {
		log.Println("Error running Nmap")
		log.Println(err.Error())
		return nil, err
	}
	output := string(outputBytes)

	return &output, nil
}

func isValidHost(address string) bool {
	addrs, err := net.LookupHost(address)
	if err != nil {
		return false
	}

	if len(addrs) == 0 {
		return false
	}

	return true
}

func parsePorts(output string) []string {
	output = strings.ReplaceAll(output, "\r\n", "\n")
	lines := strings.Split(output, "\n")

	ports := []string{}
	for _, line := range lines {
		// The only lines of Nmap's output that starts with a number is when it's
		// listing an open port
		if len(line) > 0 && unicode.IsNumber(rune(line[0])) {
			port := strings.Split(line, "/")[0]
			ports = append(ports, port)
		}
	}

	return ports
}
