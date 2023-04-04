package utils

import (
	"bufio"
	"go.uber.org/zap"
	"os"
)

func ReadFuzzFile(logger *zap.SugaredLogger, dictionaryPath string) []string {
	var words []string
	readFile, err := os.Open(dictionaryPath)

	if err != nil {
		logger.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		words = append(words, fileScanner.Text())
	}

	defer readFile.Close()

	return words
}
