package csv2json

import (
	"encoding/json"
	"io"
	"log/slog"
)

func WriteJSONFile(jsonOutput io.Writer, csvContent <-chan map[string]string, done chan<- bool) error {
	defer func() {
		done <- true
		close(done)
	}()

	counter := 0

	for csv := range csvContent {
		slog.Info("write record", "record", csv)

		if counter == 0 {
			if _, err := jsonOutput.Write([]byte("[")); err != nil {
				slog.Error("Error writing to file", "error", err)
				return err
			}
		}

		if counter > 0 {
			if _, err := jsonOutput.Write([]byte(",")); err != nil {
				slog.Error("Error writing to file", "error", err)
				return err
			}
		}

		counter++

		jsonData, err := json.Marshal(csv)
		if err != nil {
			slog.Error("Error converting to JSON", "error", err)
			continue
		}

		if _, err := jsonOutput.Write(jsonData); err != nil {
			slog.Error("Error writing to file", "error", err)
		}

	}

	if _, err := jsonOutput.Write([]byte("]")); err != nil {
		slog.Error("Error writing to file", "error", err)
	}

	return nil
}
