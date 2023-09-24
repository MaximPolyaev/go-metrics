package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/hash"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type HTTPClient struct {
	client  http.Client
	baseURL string
	hashKey *string
}

const updatesAction = "/updates/"

func NewHTTPClient(baseURL string, hashKey *string) *HTTPClient {
	return &HTTPClient{
		client:  http.Client{},
		baseURL: baseURL,
		hashKey: hashKey,
	}
}

func (c *HTTPClient) UpdateMetrics(mSlice []metric.Metric) error {
	body, err := json.Marshal(mSlice)
	if err != nil {
		return err
	}

	req, err := c.newUpdateReq(updatesAction, body)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		defer func() { _ = resp.Body.Close() }()

		errorMsg := "not update metrics, err: "

		if err != nil {
			errorMsg += err.Error()
		} else {
			errorMsg += string(body)
		}

		return errors.New(errorMsg)
	}

	return nil
}

func (c *HTTPClient) newUpdateReq(url string, body []byte) (*http.Request, error) {
	var buf bytes.Buffer

	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(body); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Encoding", "gzip")

	if c.hashKey != nil {
		req.Header.Add("HashSHA256", hash.Encode(buf.Bytes(), *c.hashKey))
	}

	return req, nil
}
