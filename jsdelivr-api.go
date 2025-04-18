package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StatsResponse represents the JSON response from JSDelivr API
type StatsResponse struct {
	Rank     int                    `json:"rank"`
	TypeRank int                    `json:"typeRank"`
	Total    int                    `json:"total"`
	Versions map[string]VersionData `json:"versions"`
}

// VersionData stores the download data for a specific version
type VersionData struct {
	Total int            `json:"total"`
	Dates map[string]int `json:"dates"`
}

// PackageStats stores the stats for a package
type PackageStats struct {
	Name      string
	Version   string
	Downloads map[string]int
}

// fetchJSDelivrStats retrieves download stats for a package from JSDelivr API
func fetchJSDelivrStats(packageName string) ([]*PackageStats, error) {
	url := fmt.Sprintf("https://data.jsdelivr.com/v1/package/npm/%s/stats", packageName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var statsResp StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statsResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var allVersionStats []*PackageStats

	// Iterate through all versions
	for version, data := range statsResp.Versions {
		stats := &PackageStats{
			Name:      packageName,
			Version:   version,
			Downloads: data.Dates,
		}
		allVersionStats = append(allVersionStats, stats)
	}

	return allVersionStats, nil
}

func collectJSDelivrStats() error {
	packages := []string{"@astrouxds/react", "@astrouxds/astro-web-components"}
	csvWriter := NewCSVWriter("jsdelivr_stats.csv")

	// Write header
	if err := csvWriter.WriteHeader(); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	for _, pkg := range packages {
		stats, err := fetchJSDelivrStats(pkg)
		if err != nil {
			fmt.Printf("Error fetching stats for %s: %v\n", pkg, err)
			continue
		}

		for _, stat := range stats {
			for date, downloads := range stat.Downloads {
				if err := csvWriter.AppendData(date, stat.Name, stat.Version, downloads); err != nil {
					return fmt.Errorf("error writing data: %v", err)
				}
			}
		}
	}
	return nil
}
