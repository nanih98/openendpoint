package httpclient

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Response struct {
	StatusCode   int
	ResponseText string
}

func requester(url string, client *http.Client) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := client.Do(req)

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
