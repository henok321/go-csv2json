package csv2json

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
)

func ReadCSVFile(csvInput io.Reader, csvContent chan<- map[string]string, bufferSize int) error {
	defer close(csvContent)

	records := make(chan []string, bufferSize)

	if err := readLines(csvInput, records); err != nil {
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

func readLines(csvInput io.Reader, records chan<- []string) error {
	defer close(records)

	csvReader := csv.NewReader(bufio.NewReader(csvInput))

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
