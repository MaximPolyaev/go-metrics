// Package httpclient for communicate metrics server
package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/hash"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type HTTPClient struct {
	client  http.Client
	baseURL string
	// hashKey - Hash key for communicate server
	hashKey string
	encoder Encoder
	localIp string
}

type Encoder interface {
	Encode(data []byte) ([]byte, error)
}

// updatesAction - action for push metrics
const updatesAction = "/updates/"

// NewHTTPClient - make new http client
func NewHTTPClient(baseURL string, hashKey string) *HTTPClient {
	localIp := getLocalIP()
	return &HTTPClient{
		client:  http.Client{},
		baseURL: baseURL,
		hashKey: hashKey,
		localIp: localIp,
	}
}

func (c *HTTPClient) WithCryptoEncoder(encoder Encoder) {
	c.encoder = encoder
}

// UpdateMetrics - push metrics to server for update
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
		readBody, readBodyErr := io.ReadAll(resp.Body)

		closeErr := resp.Body.Close()
		if closeErr != nil {
			return closeErr
		}

		errorMsg := "not update metrics, err: "

		if readBodyErr != nil {
			errorMsg += readBodyErr.Error()
		} else {
			errorMsg += string(readBody)
		}

		return errors.New(errorMsg)
	}

	return nil
}

func (c *HTTPClient) newUpdateReq(url string, body []byte) (*http.Request, error) {
	if c.encoder != nil {
		var encodeErr error
		body, encodeErr = c.encoder.Encode(body)

		if encodeErr != nil {
			return nil, encodeErr
		}
	}

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

	if c.localIp != "" {
		req.Header.Add("X-Real-IP", c.localIp)
	}

	if c.hashKey != "" {
		encodedHash, err := hash.Encode(buf.Bytes(), c.hashKey)
		if err != nil {
			return nil, err
		}

		req.Header.Add("HashSHA256", encodedHash)
	}

	return req, nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
