package igcparser

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/AurelienS/cigare/pkg/model"
)

func Parse() (model.Flight, error) {
	cmd := exec.Command(
		"lib/executables/goigc",
		"parse",
		"/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc",
		"--output-format",
		"json")

	var flight model.Flight

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
