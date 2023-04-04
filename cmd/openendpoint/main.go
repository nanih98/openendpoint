package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/nanih98/openendpoint/internal/httpclient"
	"github.com/nanih98/openendpoint/internal/logging"
	"github.com/nanih98/openendpoint/internal/providers"
	"os"
	"time"
)

func main() {
	// Argparser
	start := time.Now()
	parser := argparse.NewParser("openendpoint", "Scan open endpoints like cloud buckets")
	workers := parser.Int("w", "workers", &argparse.Options{Required: false, Help: "Number of workers (threads)", Default: 5})
	keywords := parser.StringList("k", "keyword", &argparse.Options{Required: true, Help: "Keyword for url mutations"})
	quickScan := parser.Flag("q", "quick-scan", &argparse.Options{Required: false, Default: false, Help: "Quick scan, do not create mutations from fuzz.txt file"})
	dictionaryPath := parser.String("f", "file", &argparse.Options{Required: true, Help: "Dictionary path"})
	nameserver := parser.String("n", "nameserver", &argparse.Options{Required: false, Help: "Custom nameserver", Default: "8.8.8.8 "})

	logLevel := parser.String("l", "log-level", &argparse.Options{Required: false, Help: "Log Level", Default: "info"})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
	}

	filename := "logs.log"
	logger := logging.FileLogger(filename, *logLevel)

	awsMutations := providers.AWSMutations(*keywords, *quickScan, logger, *dictionaryPath)

	logger.Log.Info(fmt.Sprintf("%d Mutations created", len(awsMutations)))

	httpclient.Fetch(awsMutations, *workers, *nameserver, logger)
	logger.Log.Info(fmt.Sprintf("all done in %s", time.Since(start)))
}
