package providers

import (
	"fmt"
	"github.com/nanih98/openendpoint/internal/logging"
	"github.com/nanih98/openendpoint/internal/utils"
)

func AWSMutations(keywords []string, quickScan bool, logger *logging.CustomLogger, dictionaryPath string) []string {
	var mutations []string

	if quickScan {
		for _, keyword := range keywords {
			mutations = append(mutations, fmt.Sprintf("https://%s.%s", keyword, S3_URL))
		}
		return mutations
	}

	// If quickScan not selected, then create mutatiosn using your keywords and fuzz.txt file or your custom dictionary
	words := utils.ReadFuzzFile(logger, dictionaryPath)

	for _, word := range words {
		for _, keyword := range keywords {
			// Appends
			mutations = append(mutations, fmt.Sprintf("https://%s%s.%s", word, keyword, S3_URL))
			mutations = append(mutations, fmt.Sprintf("https://%s.%s.%s", word, keyword, S3_URL))
			mutations = append(mutations, fmt.Sprintf("https://%s-%s.%s", word, keyword, S3_URL))

			// Prepends
			mutations = append(mutations, fmt.Sprintf("https://%s%s.%s", keyword, word, S3_URL))
			mutations = append(mutations, fmt.Sprintf("https://%s.%s.%s", keyword, word, S3_URL))
			mutations = append(mutations, fmt.Sprintf("https://%s-%s.%s", keyword, word, S3_URL))
		}
	}

	return mutations
}
