package csv2json

import (
	"encoding/json"
	"log/slog"
	"os"
)

func WriteJSONFile(path string, csvContent <-chan map[string]string, done chan<- bool) error {
	defer func() {
		done <- true
		close(done)
	}()

	file, err := os.Create(path)
	if err != nil {
		slog.Error("Error creating file", "error", err)
		return err
	}

	defer file.Close()

	first := true

	for csv := range csvContent {
		slog.Info("csv record", "record", csv)

		if first {
			if _, err := file.Write([]byte("[")); err != nil {
				slog.Error("Error writing to file", "error", err)
				return err
			}
			first = false
			continue
		}

		jsonData, err := json.Marshal(csv)
		if err != nil {
			slog.Error("Error converting to JSON", "error", err)
			continue
		}

		if _, err := file.Write(jsonData); err != nil {
			slog.Error("Error writing to file", "error", err)
		}

		if _, err := file.Write([]byte(",")); err != nil {
			slog.Error("Error writing to file", "error", err)
			return err
		}
	}

	if _, err := file.Write([]byte("]")); err != nil {
		slog.Error("Error writing to file", "error", err)
	}

	return nil
}
