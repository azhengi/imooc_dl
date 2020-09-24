package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var NormalClientHeader = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36",
}

func Get(u string, params map[string]string) (io.ReadCloser, error) {
	c := http.Client{Timeout: time.Second * time.Duration(30)}
	request, err := http.NewRequest(http.MethodGet, u, nil)
	values := request.URL.Query()
	for k, v := range params {
		values.Add(k, v)
	}
	request.URL.RawQuery = values.Encode()

	for k, v := range NormalClientHeader {
		request.Header.Add(k, v)
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func GetSegment(url string) ([]byte, error) {
	_, err := Get(url, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func DecryptResponseInfo() {
}

func Parser(body []byte, v interface{}) error {
	decode := json.NewDecoder(bytes.NewBuffer(body))
	return decode.Decode(v)
}
