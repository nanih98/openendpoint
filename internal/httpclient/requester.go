package httpclient

import (
	"io"
	"net/http"
)

type Response struct {
	StatusCode   int
	ResponseText string
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
