package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const requestTimeout = 30 * time.Second

func requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

func printResponse(raw []byte) error {
	if jsonFlag {
		fmt.Println(string(raw))
		return nil
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		fmt.Println(string(raw))
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

	fmt.Println(string(out))
	return nil
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func runAPI(method, path string, query url.Values, body any) error {
	ctx, cancel := requestContext()
	defer cancel()

	raw, _, err := api.DoJSON(ctx, method, path, query, body, nil, waitOn429Flag)
	if err != nil {
		return err
	}
	return printResponse(raw)
}
