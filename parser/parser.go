package parser

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/AurelienS/cigare/flight"
)

func Parse() (flight.Flight, error) {
	cmd := exec.Command(
		"./goigc",
		"parse",
		"/mnt/c/Users/TheGosu/Desktop/record.igc",
		"--output-format",
		"json")

	var flight flight.Flight

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command to goigc:", err)
		return flight, err
	}

	err = json.Unmarshal(output, &flight)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return flight, err
	}
	// flight.Points = flight.Points[:600]
	return flight, nil
}
