package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type HTTPClient struct {
	client  http.Client
	baseURL string
}

const (
	updateAction = "/update/"
)

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		client:  http.Client{},
		baseURL: baseURL,
	}
}

func (c *HTTPClient) UpdateMetrics(metrics []metric.Metric) error {
	for _, mm := range metrics {
		if err := c.updateMetric(&mm); err != nil {
			return err
		}
	}

	return nil
}

func (c *HTTPClient) updateMetric(mm *metric.Metric) error {
	body, err := json.Marshal(mm)
	if err != nil {
		return err
	}

	req, err := c.newUpdateReq(updateAction, body)
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

		errorMsg := "not update metricservice " + mm.ID + ", err: "

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

	return req, nil
}
