package main

import "fmt"

func main() {
	// Collect JSDelivr stats
	if err := collectJSDelivrStats(); err != nil {
		fmt.Printf("Error collecting JSDelivr stats: %v\n", err)
	}

	// Collect NPM stats
	if err := collectNPMStats(); err != nil {
		fmt.Printf("Error collecting NPM stats: %v\n", err)
	}
}
