package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/JAGBytes/domain-security-cli/internal/models"
)

const apiBaseURL = "https://api.ssllabs.com/api/v3"

type SSLabsClient struct {
	httpClient *http.Client
}

func NewSSLabsClient() *SSLabsClient {
	return &SSLabsClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *SSLabsClient) Analyze(domain string, startNew bool) (*models.Host, error) {
	params := url.Values{}
	params.Set("host", domain)
	params.Set("all", "done")
	if startNew {
		params.Set("startNew", "on")
	}

	apiURL := fmt.Sprintf("%s/analyze?%s", apiBaseURL, params.Encode())
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limit exceeded, wait before retrying")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var host models.Host
	if err := json.Unmarshal(body, &host); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &host, nil
}
func (c *SSLabsClient) Poll(domain string, maxWait time.Duration) (*models.Host, error) {
	start := time.Now()

	for {
		host, err := c.Analyze(domain, false)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Current status: %s\n", host.Status)

		if host.Status == "READY" || host.Status == "ERROR" {
			return host, nil
		}

		if time.Since(start) > maxWait {
			return nil, fmt.Errorf("timeout after %v (last status: %s)", maxWait, host.Status)
		}

		time.Sleep(10 * time.Second)
	}
}
