package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"molt/internal/moltbook"
)

const requestTimeout = 30 * time.Second

var (
	suggestedActionRegex = regexp.MustCompile(`^([A-Z]+)\s+([^\s]+)\s*(?:—\s*(.*))?$`)
	hintTranslationRegex = regexp.MustCompile(`Send a ([A-Z]+) request to (/[^\s]+)`)
)

func requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

func printResponse(raw []byte) error {
	if jsonFlag {
		// Ensure raw JSON output is sanitized to remove invalid control characters
		// that can break downstream JSON parsers.
		sanitized := make([]byte, 0, len(raw))
		for _, b := range raw {
			if (b >= 32) || (b == 9) || (b == 10) || (b == 13) {
				sanitized = append(sanitized, b)
			} else {
				// Replace invalid control chars with a space or a safe escape
				sanitized = append(sanitized, ' ')
			}
		}
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", sanitized); err != nil {
			return err
		}
		return nil
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", raw); err != nil {
			return err
		}
		return nil
	}

	return printValue(parsed)
}

func printValue(v any) error {
	var (
		out []byte
		err error
	)

	if jsonFlag {
		out, err = json.Marshal(v)
	} else {
		out, err = json.MarshalIndent(v, "", "  ")
	}
	if err != nil {
		return err
	}
 
	if _, err := fmt.Fprintf(os.Stdout, "%s\n", out); err != nil {
		return err
	}
	return nil
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
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

func translateAPIPath(method, path string) string {
	for _, mapping := range apiToCLI {
		if mapping.method == method && mapping.regex.MatchString(path) {
			return mapping.regex.ReplaceAllString(path, mapping.cli)
		}
	}
	return method + " " + path
}

func translateSuggestedAction(action string) string {
	matches := suggestedActionRegex.FindStringSubmatch(action)
	if len(matches) > 2 {
		method := matches[1]
		path := matches[2]
		desc := ""
		if len(matches) > 3 {
			desc = strings.TrimSpace(matches[3])
		}

		cliCmd := translateAPIPath(method, path)
		if desc != "" {
			return cliCmd + "  // " + desc
		}
		return cliCmd
	}
	return action
}

// translateHint tries to extract "Send a <METHOD> request to <PATH>" and translate it to CLI
func translateHint(hint string) string {
	matches := hintTranslationRegex.FindStringSubmatch(hint)
	if len(matches) == 3 {
		method := matches[1]
		path := matches[2]
		cliCmd := translateAPIPath(method, path)
		if cliCmd != method+" "+path {
			return strings.Replace(hint, matches[0], "Use `"+cliCmd+"`", 1)
		}
	}
	return hint
}

func formatAPIError(err error, raw []byte) error {
	if err == nil || jsonFlag {
		return err
	}

	var errResp moltbook.ErrorResponse
	errRespParsed := json.Unmarshal(raw, &errResp) == nil && errResp.Error != ""

	var genericErr map[string]any
	genericErrParsed := json.Unmarshal(raw, &genericErr) == nil

	if errRespParsed {
		errMsg := errResp.Error

		// Some endpoints return message as string (like Unauthorized message)
		if genericErrParsed {
			if msg, ok := genericErr["message"]; ok {
				if msgStr, ok2 := msg.(string); ok2 {
					errMsg = fmt.Sprintf("%s (%s)", errResp.Error, msgStr)
				}
			}
		}

		if errResp.RetryAfter > 0 {
			errMsg += fmt.Sprintf(" (retry after %d seconds. Consider using --wait-on-429 to automatically wait and retry)", errResp.RetryAfter)
		} else if errResp.RetryAfterSecs > 0 {
			errMsg += fmt.Sprintf(" (retry after %d seconds. Consider using --wait-on-429 to automatically wait and retry)", errResp.RetryAfterSecs)
		} else if errResp.RetryAfterMins > 0 {
			errMsg += fmt.Sprintf(" (retry after %d minutes. Consider using --wait-on-429 to automatically wait and retry)", errResp.RetryAfterMins)
		}

		if errResp.Hint != "" {
			translatedCmd := translateHint(errResp.Hint)
			if translatedCmd != errResp.Hint {
				return fmt.Errorf("✗ Command failed\n錯在哪: %s\n正確用法: molt <command> --help\n可嘗試: %s", errMsg, translatedCmd)
			}
			return fmt.Errorf("✗ Command failed\n錯在哪: %s\n正確用法: molt <command> --help\n提示: %s", errMsg, errResp.Hint)
		}
		return fmt.Errorf("✗ Command failed\n錯在哪: %s\n正確用法: molt <command> --help", errMsg)
	}

	// For validation errors like 400 Bad Request which return message as array
	if genericErrParsed {
		if msg, ok := genericErr["message"]; ok {
			return fmt.Errorf("✗ Command failed\n錯在哪: %v\n正確用法: molt <command> --help", msg)
		}
		if msg, ok := genericErr["error"]; ok {
			return fmt.Errorf("✗ Command failed\n錯在哪: %v\n正確用法: molt <command> --help", msg)
		}
	}

	return err
}

func runAPI(method, path string, query url.Values, body any) error {
	ctx, cancel := requestContext()
	defer cancel()

	raw, _, err := api.DoJSON(ctx, method, path, query, body, nil, waitOn429Flag)
	if err != nil {
		if raw != nil {
			return formatAPIError(err, raw)
		}
		return err
	}
	return printResponse(raw)
}

var uuidRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func cleanAndValidateUUID(id string) (string, error) {
	original := id
	cleaned := strings.TrimSpace(id)
	if len(cleaned) > 0 {
		runes := []rune(cleaned)
		lastRune := runes[len(runes)-1]
		if unicode.IsPunct(lastRune) {
			cleaned = string(runes[:len(runes)-1])
		}
	}

	if uuidRegex.MatchString(cleaned) {
		return cleaned, nil
	}

	return cleaned, fmt.Errorf("invalid UUID format. original: %q, cleaned: %q", original, cleaned)
}

// runAPIAndPrint executes the API, parses JSON into `out` interface, and calls formatFunc to print text if not using --json.
func runAPIAndPrint(method, path string, query url.Values, body any, out any, formatFunc func() error) error {
	ctx, cancel := requestContext()
	defer cancel()

	raw, _, err := api.DoJSON(ctx, method, path, query, body, out, waitOn429Flag)
	if err != nil {
		if raw != nil {
			return formatAPIError(err, raw)
		}
		return err
	}

	if jsonFlag {
		return printResponse(raw)
	}

	if formatFunc != nil {
		return formatFunc()
	}

	return printResponse(raw)
}
