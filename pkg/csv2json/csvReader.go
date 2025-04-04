package csv2json

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"os"
)

func ReadCSVFile(path string, csvContent chan<- map[string]string) error {
	defer close(csvContent)

	records := make(chan []string, 10)

	if err := readLines(path, records); err != nil {
		return err
	}

	first := true
	var headers []string

	for record := range records {
		if first {
			headers = record
			first = false
			continue
		}

		parsedLine, err := parseLine(record, headers)
		if err != nil {
			slog.Error("Error parsing record", "record", record, "error", err)
		}

		csvContent <- parsedLine
	}

	return nil
}

func parseLine(record, headers []string) (map[string]string, error) {
	if len(record) != len(headers) {
		return nil, errors.New("record length does not match headers length")
	}

	mappedRecord := make(map[string]string)

	for i, header := range headers {
		mappedRecord[header] = record[i]
	}
	return mappedRecord, nil
}

func readLines(path string, records chan<- []string) error {
	defer close(records)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))

	for {
		record, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			slog.Error("Error reading CSV file", "error", err)
		}
		records <- record
	}

	return nil
}
