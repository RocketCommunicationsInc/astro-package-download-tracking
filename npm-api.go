package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type VersionDownloads map[string]int

type NPMResponse struct {
	Downloads VersionDownloads `json:"downloads"`
}

func fetchNPMStats(packageName string) (*NPMResponse, error) {
	url := fmt.Sprintf("https://api.npmjs.org/versions/@astrouxds%%2F%s/last-week", packageName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var npmResp NPMResponse
	if err := json.NewDecoder(resp.Body).Decode(&npmResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &npmResp, nil
}

func collectNPMStats() error {
	packages := []string{"astro-web-components", "react"}
	csvWriter := NewCSVWriter("npm_stats.csv")

	// Write header
	if err := csvWriter.WriteHeader(); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	for _, pkg := range packages {
		stats, err := fetchNPMStats(pkg)
		if err != nil {
			fmt.Printf("Error fetching stats for %s: %v\n", pkg, err)
			continue
		}

		for version, downloads := range stats.Downloads {
			if err := csvWriter.AppendData(currentDate, pkg, version, downloads); err != nil {
				return fmt.Errorf("error writing data: %v", err)
			}
		}
	}
	return nil
}
