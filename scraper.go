package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

// saveToCSV saves the package stats to a CSV file
func saveToCSV(stats []*PackageStats, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Date", "Package", "Version", "Downloads"}); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write data rows
	for _, stat := range stats {
		for date, downloads := range stat.Downloads {
			if err := writer.Write([]string{
				date,
				stat.Name,
				stat.Version,
				fmt.Sprintf("%d", downloads),
			}); err != nil {
				return fmt.Errorf("error writing row: %v", err)
			}
		}
	}

	return nil
}

func main() {
	// List of packages to track
	packages := []string{"@astrouxds/react", "@astrouxds/astro-web-components"}

	var allStats []*PackageStats

	for _, pkg := range packages {
		stats, err := fetchJSDelivrStats(pkg)
		if err != nil {
			fmt.Printf("Error fetching stats for %s: %v\n", pkg, err)
			continue
		}

		allStats = append(allStats, stats...)
		for _, stat := range stats {
			fmt.Printf("Fetched stats for %s: Version=%s, Downloads=%v\n",
				pkg, stat.Version, stat.Downloads)
		}
	}

	// Save data to CSV
	csvFilename := "jsdelivr_stats.csv"
	if err := saveToCSV(allStats, csvFilename); err != nil {
		fmt.Printf("Error saving to CSV: %v\n", err)
		return
	}

	fmt.Printf("Data saved to %s\n", csvFilename)
}
