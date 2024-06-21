package attacker

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/quay/zlog"
	"gopkg.in/yaml.v2"
)

//go:embed assets/*
var assets embed.FS

// Questions list struct
type QuestionList struct {
	Questions []string `yaml:"questions"`
}

// Feedback struct
type Feedback struct {
	ConversationID string `yaml:"conversation_id"`
	UserFeedback   string `yaml:"user_feedback"`
	UserQuestion   string `yaml:"user_question"`
	LLMResponse    string `yaml:"llm_response"`
}

// Feedbacks list struct
type FeedbackList struct {
	Feedbacks []Feedback `yaml:"feedbacks"`
}

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

// CreateQueryRequests returns the list of requests to perform POST operation on query endpoint.
func CreateQueryRequests(ctx context.Context, hitSize int, host, token string, withCache bool) []map[string]interface{} {
	url, headers := getRequestCommons(ctx, "/v1/query", host, token)
	var requests []map[string]interface{}

	data, err := assets.ReadFile("assets/questions.yaml")
	if err != nil {
		fmt.Errorf("error: %v", err)
	}

	var questionList QuestionList
	err = yaml.Unmarshal(data, &questionList)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questionList.Questions), func(i, j int) {
		questionList.Questions[i], questionList.Questions[j] = questionList.Questions[j], questionList.Questions[i]
	})

	body := make(map[string]string)
	if withCache {
		body["conversation_id"] = "00000000-0000-0000-0000-000000000000"
		zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for POST operation on /v1/query with cache")
	} else {
		zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for POST operation on /v1/query")
	}

	for idx := 0; idx < hitSize; idx++ {
		body["query"] = questionList.Questions[idx%len(questionList.Questions)]
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			fmt.Errorf("Error marshaling body: %v", err)
		}
		requests = append(requests, map[string]interface{}{
			"method": http.MethodPost,
			"url":    url,
			"header": headers,
			"body":   bodyBytes,
		})
	}
	return requests
}

// CreateReadinessRequests returns the list of requests to perform GET operation on readiness endpoint.
func CreateReadinessRequests(ctx context.Context, hitSize int, host string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for GET operation on /readiness")
	url, headers := getRequestCommons(ctx, "/readiness", host, "")
	var requests []map[string]interface{}

	for idx := 0; idx < hitSize; idx++ {
		requests = append(requests, map[string]interface{}{
			"method": http.MethodGet,
			"url":    url,
			"header": headers,
		})
	}
	return requests
}

// CreateLivenessRequests returns the list of requests to perform GET operation on liveness endpoint.
func CreateLivenessRequests(ctx context.Context, hitSize int, host string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for GET operation on /liveness")
	url, headers := getRequestCommons(ctx, "/liveness", host, "")
	var requests []map[string]interface{}

	for idx := 0; idx < hitSize; idx++ {
		requests = append(requests, map[string]interface{}{
			"method": http.MethodGet,
			"url":    url,
			"header": headers,
		})
	}
	return requests
}

// CreateMetricsRequests returns the list of requests to perform GET operation on metrics endpoint.
func CreateMetricsRequests(ctx context.Context, hitSize int, host, token string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for GET operation on /metrics")
	url, headers := getRequestCommons(ctx, "/metrics", host, token)
	var requests []map[string]interface{}

	for idx := 0; idx < hitSize; idx++ {
		requests = append(requests, map[string]interface{}{
			"method": http.MethodGet,
			"url":    url,
			"header": headers,
		})
	}
	return requests
}

// CreateAuthorizedRequests returns the list of requests to perform POST operation on authorized endpoint.
func CreateAuthorizedRequests(ctx context.Context, hitSize int, host, token string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for POST operation on /authorized")
	url, headers := getRequestCommons(ctx, "/authorized", host, token)
	var requests []map[string]interface{}

	for idx := 0; idx < hitSize; idx++ {
		requests = append(requests, map[string]interface{}{
			"method": http.MethodPost,
			"url":    url,
			"header": headers,
		})
	}
	return requests
}

// CreateGetFeedbackStatusRequests returns the list of requests to perform GET operation on feedback status endpoint.
func CreateGetFeedbackStatusRequests(ctx context.Context, hitSize int, host, token string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for GET operation on /v1/feedback/status")
	url, headers := getRequestCommons(ctx, "/v1/feedback/status", host, token)
	var requests []map[string]interface{}

	for idx := 0; idx < hitSize; idx++ {
		requests = append(requests, map[string]interface{}{
			"method": http.MethodGet,
			"url":    url,
			"header": headers,
		})
	}
	return requests
}

// CreateFeedbackRequests returns the list of requests to perform POST operation on feedback endpoint.
func CreateFeedbackRequests(ctx context.Context, hitSize int, host, token string) []map[string]interface{} {
	zlog.Info(ctx).Int("number of requests", hitSize).Msg("preparing requests for POST operation on /v1/feedback")
	url, headers := getRequestCommons(ctx, "/v1/feedback", host, token)
	var requests []map[string]interface{}

	data, err := assets.ReadFile("assets/feedbacks.yaml")
	if err != nil {
		fmt.Errorf("error: %v", err)
	}

	var feedbackList FeedbackList
	err = yaml.Unmarshal(data, &feedbackList)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(feedbackList.Feedbacks), func(i, j int) {
		feedbackList.Feedbacks[i], feedbackList.Feedbacks[j] = feedbackList.Feedbacks[j], feedbackList.Feedbacks[i]
	})

	for idx := 0; idx < hitSize; idx++ {
		feedback := feedbackList.Feedbacks[idx%len(feedbackList.Feedbacks)]
		body := map[string]string{
			"conversation_id": feedback.ConversationID,
			"user_feedback":   feedback.UserFeedback,
			"user_question":   feedback.UserQuestion,
			"llm_response":    feedback.LLMResponse,
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			fmt.Errorf("Error marshaling body: %v", err)
		}
		requests = append(requests, map[string]interface{}{
			"method": http.MethodPost,
			"url":    url,
			"header": headers,
			"body":   bodyBytes,
		})
	}
	return requests
}
