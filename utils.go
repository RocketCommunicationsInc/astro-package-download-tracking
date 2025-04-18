package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// CSVWriter handles writing stats to CSV files
type CSVWriter struct {
	filename string
}

func NewCSVWriter(filename string) *CSVWriter {
	return &CSVWriter{filename: filename}
}

func (w *CSVWriter) WriteHeader() error {
	file, err := os.Create(w.filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Date", "Package", "Version", "Downloads"}); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}
	return nil
}

func (w *CSVWriter) AppendData(date, packageName, version string, downloads int) error {
	file, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		date,
		packageName,
		version,
		fmt.Sprintf("%d", downloads),
	}); err != nil {
		return fmt.Errorf("error writing row: %v", err)
	}
	return nil
}
