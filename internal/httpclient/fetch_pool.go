package httpclient

import (
	"fmt"
	"github.com/nanih98/openendpoint/internal/logging"
	"github.com/nanih98/openendpoint/internal/utils"
	"sync"
)

func Fetch(urls []string, workers int, nameserver string, logger *logging.CustomLogger) {
	var errs []error

	workQueue := make(chan string, len(urls))

	wg := sync.WaitGroup{}
	wg.Add(workers)

	// HTTP CLIENT
	client := HTTPClient(nameserver)

	for i := 0; i < workers; i++ {
		worker := i
		go func(worker int, workQueue chan string) {
			for uri := range workQueue {
				//start := time.Now()
				response, err := requester(uri, client)

				if err != nil {
					errs = append(errs, err)
				}

				if response.StatusCode == 200 {
					logger.Log.Info(fmt.Sprintf("Opened bucket %s", uri))
					// List content
					utils.ListBucketContents(response.ResponseText, uri)
				} else if response.StatusCode == 403 {
					logger.Log.Debug(fmt.Sprintf("Protected bucket %s", uri))
				}
			}
			wg.Done()
		}(worker, workQueue)
	}

	go func() {
		for _, url := range urls {
			workQueue <- url
		}
		close(workQueue)
	}()
	wg.Wait()

	fmt.Println(errs)
}
