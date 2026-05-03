package moltbook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"molt/internal/logger"
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
	logger.SetAPICalled()

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
		logger.SetAPIFailed()
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
				fmt.Fprintf(os.Stderr, "Rate limited. Waiting %d seconds...\n", retryAfter)
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
			errMsg += " (Consider using --wait-on-429 to automatically wait and retry)"
			return rawResp, resp.StatusCode, fmt.Errorf("rate limit: %s", errMsg)
		}
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		logger.SetAPIFailed()
		var errResp ErrorResponse
		if err := json.Unmarshal(rawResp, &errResp); err == nil && errResp.Error != "" {
			errMsg := errResp.Error
			if errResp.Hint != "" {
				errMsg += " (hint: " + TranslateHint(errResp.Hint) + ")"
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


var apiToCLI = []struct {
	method string
	regex  *regexp.Regexp
	cli    string
}{
	// Standard hints without /api/v1 prefix
	{"POST", regexp.MustCompile(`^/agents/dm/requests/([^/]+)/approve$`), "molt dm approve $1"},
	{"POST", regexp.MustCompile(`^/agents/dm/requests/([^/]+)/reject$`), "molt dm reject $1"},
	{"POST", regexp.MustCompile(`^/agents/dm/conversations/([^/]+)/send$`), "molt dm send $1"},
	{"POST", regexp.MustCompile(`^/agents/dm/request$`), "molt dm request"},
	{"GET", regexp.MustCompile(`^/agents/dm/conversations/([^/]+)$`), "molt dm read $1"},
	{"GET", regexp.MustCompile(`^/agents/dm/conversations$`), "molt dm conversations"},
	{"GET", regexp.MustCompile(`^/agents/dm/requests$`), "molt dm requests"},
	{"GET", regexp.MustCompile(`^/agents/dm/check$`), "molt dm check"},

	{"POST", regexp.MustCompile(`^/posts/([^/]+)/comments$`), "molt comment add $1"},
	{"POST", regexp.MustCompile(`^/posts/([^/]+)/upvote$`), "molt vote post-up $1"},
	{"POST", regexp.MustCompile(`^/posts/([^/]+)/downvote$`), "molt vote post-down $1"},
	{"POST", regexp.MustCompile(`^/comments/([^/]+)/upvote$`), "molt vote comment-up $1"},

	{"POST", regexp.MustCompile(`^/posts$`), "molt post create"},
	{"GET", regexp.MustCompile(`^/posts/([^/]+)$`), "molt post get $1"},
	{"DELETE", regexp.MustCompile(`^/posts/([^/]+)$`), "molt post delete $1"},

	{"GET", regexp.MustCompile(`^/search$`), "molt search"},

	// Home dashboard actions with /api/v1 prefix
	{"POST", regexp.MustCompile(`^/api/v1/notifications/read-by-post/([^/\?]+)(?:\?.*)?$`), "molt notifications read-post $1"},
	{"POST", regexp.MustCompile(`^/api/v1/notifications/read-all$`), "molt notifications read-all"},
	{"GET", regexp.MustCompile(`^/api/v1/posts/([^/\?]+)/comments(?:\?.*)?$`), "molt post comments $1"},
	{"POST", regexp.MustCompile(`^/api/v1/posts/([^/\?]+)/comments(?:\?.*)?$`), "molt comment add $1"},
	{"GET", regexp.MustCompile(`^/api/v1/posts/([^/\?]+)(?:\?.*)?$`), "molt post get $1"},
	{"GET", regexp.MustCompile(`^/api/v1/feed(?:\?filter=following.*)?$`), "molt feed my --filter following"},
	{"GET", regexp.MustCompile(`^/api/v1/feed(?:\?.*)?$`), "molt feed my"},
	{"GET", regexp.MustCompile(`^/api/v1/posts(?:\?.*)?$`), "molt feed global"},
}

func TranslateAPIPath(method, path string) string {
	for _, mapping := range apiToCLI {
		if mapping.method == method && mapping.regex.MatchString(path) {
			return mapping.regex.ReplaceAllString(path, mapping.cli)
		}
	}
	return method + " " + path
}

func TranslateSuggestedAction(action string) string {
	regex := regexp.MustCompile(`^([A-Z]+)\s+([^\s]+)\s*(?:—\s*(.*))?$`)
	matches := regex.FindStringSubmatch(action)
	if len(matches) > 2 {
		method := matches[1]
		path := matches[2]
		desc := ""
		if len(matches) > 3 {
			desc = strings.TrimSpace(matches[3])
		}

		cliCmd := TranslateAPIPath(method, path)
		if desc != "" {
			return cliCmd + "  // " + desc
		}
		return cliCmd
	}
	return action
}

// TranslateHint tries to extract "Send a <METHOD> request to <PATH>" and translate it to CLI
func TranslateHint(hint string) string {
	hintRegex := regexp.MustCompile(`Send a ([A-Z]+) request to (/[^\s]+)`)
	matches := hintRegex.FindStringSubmatch(hint)
	if len(matches) == 3 {
		method := matches[1]
		path := matches[2]
		cliCmd := TranslateAPIPath(method, path)
		if cliCmd != method+" "+path {
			return strings.Replace(hint, matches[0], "Use `"+cliCmd+"`", 1)
		}
	}
	return hint
}
