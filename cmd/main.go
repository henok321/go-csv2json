package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/henok321/go-csv2json/pkg/csv2json"
)

func main() {
	csvFile := flag.String("csvInput", "", "csv input file")
	jsonFile := flag.String("jsonOutput", "", "json output file")

	flag.Parse()

	csvContent := make(chan map[string]string, 10)
	done := make(chan bool, 1)

	go func() {
		if err := csv2json.ReadCSVFile(*csvFile, csvContent); err != nil {
			slog.Error("error reading csv file", "error", err)
		}
	}()

	go func() {
		if err := csv2json.WriteJSONFile(*jsonFile, csvContent, done); err != nil {
			slog.Error("error writing json file", "error", err)
			os.Exit(1)
		}
	}()

	<-done
}
