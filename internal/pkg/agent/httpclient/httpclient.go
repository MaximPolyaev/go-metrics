package httpclient

import (
	"errors"
	"fmt"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/metric"
	"io"
	"net/http"
)

type HTTPClient interface {
	UpdateMetrics(stats *metric.Stats) error
}

type httpClient struct {
	client  http.Client
	baseURL string
}

const (
	updateAction = "/update/"
)

func NewHTTPClient(baseURL string) HTTPClient {
	return &httpClient{
		client:  http.Client{},
		baseURL: baseURL,
	}
}

func (c *httpClient) UpdateMetrics(stats *metric.Stats) error {
	for k, v := range stats.GetCounterMap() {
		if err := c.updateCounterMetric(k, v); err != nil {
			return err
		}
	}

	for k, v := range stats.GetGaugeMap() {
		if err := c.updateGaugeMetric(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (c *httpClient) updateGaugeMetric(name string, value float64) error {
	url := c.makeUpdateURL(
		metric.GaugeType,
		name,
		fmt.Sprintf("%f", value),
	)

	return c.updateMetric(url)
}

func (c *httpClient) updateCounterMetric(name string, value int) error {
	url := c.makeUpdateURL(
		metric.CounterType,
		name,
		fmt.Sprintf("%d", value),
	)

	return c.updateMetric(url)
}

func (c *httpClient) makeUpdateURL(args ...string) string {
	url := c.baseURL + updateAction

	for _, arg := range args {
		url += arg + "/"
	}

	return url
}

func (c *httpClient) updateMetric(url string) error {
	req, err := c.makeRequest(url)
	if err != nil {
		return err
	}

	fmt.Println(url)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		defer func() { _ = resp.Body.Close() }()

		errorMsg := "not update metric " + url + ", err: "

		if err != nil {
			errorMsg += err.Error()
		} else {
			errorMsg += string(body)
		}

		return errors.New(errorMsg)
	}

	return nil
}

func (c *httpClient) makeRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/plain")

	return req, nil
}
