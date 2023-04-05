package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/nanih98/openendpoint/internal/httpclient"
	"github.com/nanih98/openendpoint/internal/logging"
	"github.com/nanih98/openendpoint/internal/providers"
	"log"
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
	logLevel := parser.Selector("l", "log-level", []string{"info", "debug"}, &argparse.Options{Required: false, Default: "info", Help: "Log level of the application"})
	//logFile := parser.File("", "log-file", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644, &argparse.Options{Required: false, Default: nil, Help: "Log file path"})
	err := parser.Parse(os.Args)

	if err != nil {
		log.Fatal(fmt.Println(parser.Usage(err)))
	}

	logger := logging.NewLogger(&logging.LoggerOptions{
		//LogFilePath: *logFile,
		LogLevel: *logLevel,
		Sugared:  true,
	})

	messages := make(chan string)
	errors := make(chan error)

	awsMutations := providers.AWSMutations(*keywords, *quickScan, logger, *dictionaryPath)

	logger.Log.Info(fmt.Sprintf("%d Mutations created", len(awsMutations)))

	go func() {
		httpclient.Fetch(awsMutations, *workers, *nameserver, logger, messages, errors)
	}()

	// Read all the application messages
	for {
		select {
		case msg := <-messages:
			logger.Log.Info(msg)
		case error := <-errors:
			// Also prints error messages
			logger.Log.Error(error)
		}
	}

	logger.Log.Info(fmt.Sprintf("all done in %s", time.Since(start)))
}
