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
	dictionaryPath := parser.String("d", "dictionary", &argparse.Options{Required: true, Help: "Dictionary path"})
	nameserver := parser.String("n", "nameserver", &argparse.Options{Required: false, Help: "Custom nameserver", Default: "8.8.8.8 "})
	logLevel := parser.Selector("l", "log-level", []string{"info", "debug"}, &argparse.Options{Required: false, Default: "info", Help: "Log level of the application"})
	//logFile := parser.File("", "log-file", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644, &argparse.Options{Required: false, Default: nil, Help: "Log file path"})
	//logFile := parser.String("f", "log-file", &argparse.Options{Required: false, Default: "", Help: "Save logs into file if set"})
	err := parser.Parse(os.Args)

	if err != nil {
		log.Fatal(fmt.Println(parser.Usage(err)))
	}

	logger := logging.NewLogger(&logging.LoggerOptions{
		LogLevel: *logLevel,
		Sugar:    true,
	})

	messages := make(chan string)
	requestError := make(chan string)
	errors := make(chan error)
	requestError := make(chan string)

	awsMutations := providers.AWSMutations(*keywords, *quickScan, logger, *dictionaryPath)

	logger.Log.Info(fmt.Sprintf("%d Mutations created", len(awsMutations)))

	go httpclient.Fetch(awsMutations, *workers, *nameserver, logger, messages, requestError, errors)

	// Read all the application messages from channels
	for {
		select {
		case msg := <-messages:
			logger.Log.Info(msg)
		case re := <-requestError:
			logger.Log.Debug(re)
		case error := <-errors:
			// Also prints error messages
			logger.Log.Error(error)
		default:
			logger.Log.Info("no work")
		}
	}

	logger.Log.Info(fmt.Sprintf("all done in %s", time.Since(start)))
}
