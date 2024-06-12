package attacker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quay/zlog"
)

// getRequestCommons returns the common inputs for each and every HTTP request.
// It returns a url and headers for the specified input.
func getRequestCommons(ctx context.Context, endpoint, host, token string) (string, http.Header) {
	url := host + endpoint
	zlog.Info(ctx).Str("endpoint", url).Msg("preparing headers")
	headers := http.Header{
		"accept":        []string{"application/json"},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
	}
	return url, headers
}

// CreateQueryRequests returns the list of requests to perform POST operation in query endpoint.
func CreateQueryRequests(ctx context.Context, hitSize int, host, token string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for POST operation in /query")
	url, headers := getRequestCommons(ctx, "/v1/query", host, token)
	var requests []map[string]interface{}
	body := map[string]string{
		"query": "write a deployment yaml for the mongodb image",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("Error marshaling body: %v", err)
	}
	for idx := 0; idx < hitSize; idx++ {
		requests = append(requests, map[string]interface{}{
			"method": http.MethodPost,
			"url":    url,
			"header": headers,
			"body":   bodyBytes,
		})
	}
	return requests
}
