package csv2json

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
)

func ReadCSVFile(csvInput io.Reader, csvDataChannel chan<- map[string]string, bufferSize int) error {
	rawRecords := make(chan []string, bufferSize)
	defer close(csvDataChannel)

	go func() {
		if err := readLines(csvInput, rawRecords); err != nil {
			slog.Error("Error reading CSV file", "error", err)
		}
	}()

	first := true
	var headers []string

	for record := range rawRecords {
		if first {
			headers = record
			first = false
			continue
		}

		parsedLine, err := parseRecord(record, headers)
		if err != nil {
			slog.Error("Error parsing record", "record", record, "error", err)
		}

		slog.Debug("parse record", "record", record)

		csvDataChannel <- parsedLine
	}

	return nil
}

func parseRecord(record, headers []string) (map[string]string, error) {
	if len(record) != len(headers) {
		return nil, errors.New("record length does not match headers length")
	}

	mappedRecord := make(map[string]string)

	for i, header := range headers {
		mappedRecord[header] = record[i]
	}
	return mappedRecord, nil
}

func readLines(csvInput io.Reader, rawRecords chan<- []string) error {
	defer close(rawRecords)

	csvReader := csv.NewReader(bufio.NewReader(csvInput))

	for {
		record, err := csvReader.Read()
		slog.Debug("read record from file", "record", record)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			slog.Error("Error reading CSV file", "error", err)
		}
		rawRecords <- record
	}

	return nil
}
