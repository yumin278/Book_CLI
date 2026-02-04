package moltbook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURL is the forced base URL with 'www' to avoid redirects
	BaseURL = "https://www.moltbook.com/api/v1"
)

// Client represents a Moltbook API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// NewClient creates a new Moltbook API client
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: apiKey,
	}
}

// DoJSON performs an HTTP request and handles JSON response
// Returns: raw response bytes, status code, error
func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, body interface{}, out interface{}, waitOn429 bool) ([]byte, int, error) {
	// Build URL
	fullURL := c.BaseURL + path
	if query != nil && len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	// Prepare request body
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle 429 rate limit
	if resp.StatusCode == http.StatusTooManyRequests {
		var errResp ErrorResponse
		if err := json.Unmarshal(rawResp, &errResp); err == nil {
			retryAfter := errResp.RetryAfter
			if retryAfter == 0 && errResp.RetryAfterSecs > 0 {
				retryAfter = errResp.RetryAfterSecs
			}
			if retryAfter == 0 && errResp.RetryAfterMins > 0 {
				retryAfter = errResp.RetryAfterMins * 60
			}

			if waitOn429 && retryAfter > 0 {
				fmt.Fprintf(io.Discard, "Rate limited. Waiting %d seconds...\n", retryAfter)
				time.Sleep(time.Duration(retryAfter) * time.Second)
				// Retry once after waiting
				return c.DoJSON(ctx, method, path, query, body, out, false)
			}

			// Format error message with retry info
			errMsg := errResp.Error
			if retryAfter > 0 {
				if retryAfter >= 60 {
					errMsg += fmt.Sprintf(" (retry after %d minutes)", retryAfter/60)
				} else {
					errMsg += fmt.Sprintf(" (retry after %d seconds)", retryAfter)
				}
			}
			if errResp.DailyRemaining >= 0 {
				errMsg += fmt.Sprintf(" [daily remaining: %d]", errResp.DailyRemaining)
			}
			return rawResp, resp.StatusCode, fmt.Errorf("rate limit: %s", errMsg)
		}
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(rawResp, &errResp); err == nil && errResp.Error != "" {
			errMsg := errResp.Error
			if errResp.Hint != "" {
				errMsg += " (hint: " + errResp.Hint + ")"
			}
			return rawResp, resp.StatusCode, fmt.Errorf("API error (%d): %s", resp.StatusCode, errMsg)
		}
		return rawResp, resp.StatusCode, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(rawResp))
	}

	// Parse successful response
	if out != nil {
		if err := json.Unmarshal(rawResp, out); err != nil {
			return rawResp, resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return rawResp, resp.StatusCode, nil
}
