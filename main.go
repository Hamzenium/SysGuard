package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runBinary(binaryPath string) {
	// Run the binary
	cmd := exec.Command(binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running binary: %s\n", binaryPath)
	err := cmd.Run() // Execute the command
	if err != nil {
		fmt.Printf("Error running binary %s: %v\n", binaryPath, err)
	}
}

func main() {
	// Path to your binary file
	binaryPath := "frontend/sysguard"

	// Run the binary only once
	runBinary(binaryPath)

	fmt.Println("Binary execution has finished.")
}
