package main

import (
	"log"
	"os/exec"
	"strings"
	"unicode"
)

func getOpenPorts(address string) []string {
	output := runNmap(address)
	return parsePorts(output)
}

func runNmap(address string) string {
	log.Println("Running Nmap for \"" + address + "\"")

	// TODO: Sanitize address to ensure it is a valid IP/host

	output, err := exec.Command("nmap", "--open", "-p0-1000", address).Output()
	if err != nil {
		log.Println("Error running Nmap")
		log.Println(err.Error())
	}

	return string(output)
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
