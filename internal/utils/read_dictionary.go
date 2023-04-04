package utils

import (
	"bufio"
	"github.com/nanih98/openendpoint/internal/logging"
	"os"
)

func ReadFuzzFile(logger logging.Logger, dictionaryPath string) []string {
	var words []string
	readFile, err := os.Open(dictionaryPath)

	if err != nil {
		logger.Log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		words = append(words, fileScanner.Text())
	}

	defer readFile.Close()

	return words
}
