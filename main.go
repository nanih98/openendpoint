package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/nanih98/openendpoint/internal/utils"
	"github.com/nanih98/openendpoint/logger"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Response struct {
	StatusCode   int
	ResponseText string
}

const (
	S3_URL = "s3.amazonaws.com"
)

func ReadFuzzFile() []string {
	var words []string
	//readFile, err := os.Open("/usr/local/share/SecLists/Discovery/Web-Content/directory-list-2.3-small.txt")
	readFile, err := os.Open("fuzz.txt")

	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		words = append(words, fileScanner.Text())
	}
	readFile.Close()

	return words
}

func Mutations(keywords []string, quickScan bool) []string {
	var mutations []string

	if quickScan {
		for _, keyword := range keywords {
			mutations = append(mutations, fmt.Sprintf("https://%s.%s", keyword, S3_URL))
		}
		return mutations
	}

	// If quickScan not selected, then create mutatiosn using your keywords and fuzz.txt file or your custom dictionary
	words := ReadFuzzFile()

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

func fetch(urls []string, workers int, nameserver string, logger *zap.SugaredLogger) {
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
				} else {
					logger.Warn(fmt.Sprintf("Not found %s %d", uri, response.StatusCode))
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

func HTTPClient(nameserver string) *http.Client {
	var (
		dnsResolverIP        = nameserver + ":53"
		dnsResolverProto     = "udp"
		dnsResolverTimeoutMs = 5000
	)

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
				}
				return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
			},
		},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	tr := &http.Transport{
		MaxIdleConns:          50,
		IdleConnTimeout:       30 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		DisableCompression:    true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DialContext:           dialContext,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 15,
	}

	return client
}

func requester(url string, client *http.Client) (Response, error) {
	resp, err := client.Get(url)

	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)

	if err != nil {
		return Response{}, err
	}

	responseString := string(responseData)

	return Response{StatusCode: resp.StatusCode, ResponseText: responseString}, nil
}

func main() {
	// Argparser
	start := time.Now()
	parser := argparse.NewParser("openendpoint", "Scan open endpoints like cloud buckets")
	workers := parser.Int("w", "workers", &argparse.Options{Required: false, Help: "Number of workers (threads)", Default: 5})
	keywords := parser.StringList("k", "keyword", &argparse.Options{Required: true, Help: "Keyword for url mutations"})
	quickScan := parser.Flag("q", "quick-scan", &argparse.Options{Required: false, Default: false, Help: "Quick scan, do not create mutations from fuzz.txt file"})
	//dictionaryPath := parser.String("f", "file", &argparse.Options{Required: false, Help: "Dictionary path"})
	nameserver := parser.String("n", "nameserver", &argparse.Options{Required: false, Help: "Custom nameserver", Default: "8.8.8.8 "})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
	}

	filename := "logs.log"
	logger := logger.FileLogger(filename)

	mutations := Mutations(*keywords, *quickScan)

	logger.Info(fmt.Sprintf("%d Mutations created", len(mutations)))

	fetch(mutations, *workers, *nameserver, logger)
	logger.Info(fmt.Sprintf("all done in %s", time.Since(start)))
}
