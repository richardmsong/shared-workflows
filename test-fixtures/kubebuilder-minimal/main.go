// Package main provides a minimal kubebuilder-style controller manager.
package main

import (
	"fmt"
	"os"
)

// Version information set via ldflags
var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	fmt.Printf("kubebuilder-minimal version %s (commit: %s)\n", version, commit)
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// In a real kubebuilder project, this would set up the controller manager
	// For this minimal example, we just demonstrate the structure
	fmt.Println("Starting controller manager...")
	return nil
}

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}
