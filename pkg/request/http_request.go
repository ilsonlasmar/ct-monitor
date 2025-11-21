package request

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpRequest struct {
	client http.Client
	url    string
}

func NewHttpRequest(client http.Client, url string) *HttpRequest {
	return &HttpRequest{
		client: client,
		url:    url,
	}
}

func (h *HttpRequest) Get() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, h.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
		return nil, err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}
