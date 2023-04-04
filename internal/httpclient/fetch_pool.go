package httpclient

import (
	"fmt"
	"github.com/nanih98/openendpoint/internal/utils"
	"go.uber.org/zap"
	"sync"
)

func Fetch(urls []string, workers int, nameserver string, logger *zap.SugaredLogger) {
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
					logger.Info(fmt.Sprintf("Opened bucket %s", uri))
					// List content
					utils.ListBucketContents(response.ResponseText, uri)
				} else if response.StatusCode == 403 {
					logger.Debug(fmt.Sprintf("Protected bucket %s", uri))
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
