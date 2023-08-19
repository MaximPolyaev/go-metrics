package httpclient

import (
	"bytes"
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

func (c *HTTPClient) UpdateMetrics(metrics []metric.Metrics) error {
	for _, mm := range metrics {
		if err := c.updateMetric(&mm); err != nil {
			return err
		}
	}

	return nil
}

func (c *HTTPClient) updateMetric(mm *metric.Metrics) error {
	body, err := json.Marshal(mm)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+updateAction, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		defer func() { _ = resp.Body.Close() }()

		errorMsg := "not update metric " + mm.ID + ", err: "

		if err != nil {
			errorMsg += err.Error()
		} else {
			errorMsg += string(body)
		}

		return errors.New(errorMsg)
	}

	return nil
}
