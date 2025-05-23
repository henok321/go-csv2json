package main

import (
	"flag"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/henok321/go-csv2json/pkg/csv2json"
)

func main() {
	csvFile := flag.String("csvInput", "", "csv input file")
	jsonFile := flag.String("jsonOutput", "", "json output file")
	bufferSize := flag.Int("bufferSize", 1, "number of buffer for parallel conversion (default 1)")

	flag.Parse()

	if *csvFile == "" || *jsonFile == "" {
		slog.Error("csvInput and jsonOutput flags are required")
		flag.Usage()
		os.Exit(1)
	}

	csvInput, err := os.Open(*csvFile)
	if err != nil {
		slog.Error("error opening csv file", "error", err)
		os.Exit(1)
	}

	jsonOutput, err := os.Create(*jsonFile)
	if err != nil {
		slog.Error("error opening json file", "error", err)
		os.Exit(1)
	}

	slog.Info("starting conversion", "csvInput", *csvFile, "jsonOutput", *jsonFile, "bufferSize", *bufferSize)

	StartConversion(csvInput, jsonOutput, *bufferSize)
}

func StartConversion(csvInput io.Reader, jsonOutput io.Writer, bufferSize int) {
	defer track("total conversion")()
	csvDataChannel := make(chan map[string]string, bufferSize)
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer track("readCSV")()
		defer wg.Done()
		if err := csv2json.ReadCSVFile(csvInput, csvDataChannel, bufferSize); err != nil {
			slog.Error("error reading csv file", "error", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer track("writeJSON")()
		defer wg.Done()
		if err := csv2json.WriteJSONFile(jsonOutput, csvDataChannel); err != nil {
			slog.Error("error writing json file", "error", err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}

func track(name string) func() {
	start := time.Now()
	return func() {
		slog.Info("execution time", "name", name, "duration", time.Since(start))
	}
}
