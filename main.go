package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/nanih98/openendpoint/logger"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

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

func Mutations(keyword string) []string {
	var mutations []string

	words := ReadFuzzFile()

	for _, word := range words {
		// Appends
		mutations = append(mutations, fmt.Sprintf("http://%s%s.%s", word, keyword, S3_URL))
		mutations = append(mutations, fmt.Sprintf("http://%s.%s.%s", word, keyword, S3_URL))
		mutations = append(mutations, fmt.Sprintf("http://%s-%s.%s", word, keyword, S3_URL))

		// Prepends
		mutations = append(mutations, fmt.Sprintf("http://%s%s.%s", keyword, word, S3_URL))
		mutations = append(mutations, fmt.Sprintf("http://%s.%s.%s", keyword, word, S3_URL))
		mutations = append(mutations, fmt.Sprintf("http://%s-%s.%s", keyword, word, S3_URL))
	}
	return mutations
}

func fetch(urls []string, workers int, nameserver string, logger *zap.Logger) {
	//var errs []error

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
				status := requester(uri, client)

				//if err != nil {
				//	errs = append(errs, err)
				//}

				if status == 200 {
					logger.Info("Opened bucket",
						zap.Int("worker", worker),
						zap.String("url", uri),
						zap.Int("status", status),
					)

				} else if status == 403 {
					logger.Warn("Protected bucket",
						zap.Int("worker", worker),
						zap.String("url", uri),
						zap.Int("status", status),
					)
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
		Timeout:   time.Second * 5,
	}

	return client
}

func requester(url string, client *http.Client) int {
	resp, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode
}

func main() {
	var fuzzFile, nameserver, keyword string
	var threads int

	flag.IntVar(&threads, "threads", 5, "Number of threads")
	flag.StringVar(&keyword, "keyword", "", "Keyword for mutation")
	flag.StringVar(&fuzzFile, "fuzzFilePath", "fuzz.txt", "Specify fuzz.txt file path")
	flag.StringVar(&nameserver, "nameserver", "8.8.8.8", "Specify custom nameserver")
	flag.Parse()

	filename := "logs.log"
	logger := logger.FileLogger(filename)

	start := time.Now()
	mutations := Mutations(keyword)
	logger.Info(fmt.Sprintf("Mutations created in (%.2fs)", time.Since(start).Seconds()))

	fetch(mutations, threads, nameserver, logger)
}
