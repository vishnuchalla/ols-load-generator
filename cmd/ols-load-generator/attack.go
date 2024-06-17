package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/quay/zlog"
	"github.com/urfave/cli/v2"
	"github.com/vishnuchalla/ols-load-generator/attacker"
)

// Command line to handle attack functionality.
var AttackCmd = &cli.Command{
	Name:        "attack",
	Description: "perform attack on ols endpoints",
	Usage:       "ols-load-test attack",
	Action:      attackAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "host",
			Usage:   "--host localhost:6060/",
			Value:   "http://localhost:6060/",
			EnvVars: []string{"OLS_TEST_HOST"},
		},
		&cli.StringFlag{
			Name:    "runid",
			Usage:   "--runid f519d9b2-aa62-44ab-9ce8-4156b712f6d2",
			Value:   uuid.New().String(),
			EnvVars: []string{"OLS_TEST_RUNID"},
		},
		&cli.StringFlag{
			Name:    "authtoken",
			Usage:   "--authtoken authtoken",
			Value:   "",
			EnvVars: []string{"OLS_TEST_AUTH_TOKEN"},
		},
		&cli.IntFlag{
			Name:    "hitsize",
			Usage:   "--hitsize 100",
			Value:   25,
			EnvVars: []string{"OLS_TEST_HIT_SIZE"},
		},
		&cli.IntFlag{
			Name:    "rps",
			Usage:   "--rps 50",
			Value:   10,
			EnvVars: []string{"OLS_TEST_RPS"},
		},
	},
}

// Type to store the test config.
type TestConfig struct {
	Rps       int    `json:"rps"`
	Host      string `json:"host"`
	HitSize   int    `json:"hitsize"`
	AuthToken string `json:"-"`
	RunID     string `json:"runid"`
}

// NewConfig creates and returns a test configuration from CLI options.
func NewConfig(c *cli.Context) *TestConfig {
	return &TestConfig{
		AuthToken: c.String("authtoken"),
		RunID:     c.String("runid"),
		Host:      c.String("host"),
		HitSize:   c.Int("hitsize"),
		Rps:       c.Int("rps"),
	}
}

// attackAction drives the attacking logic.
// It returns an error if any during the execution.
func attackAction(c *cli.Context) error {
	startTime := time.Now()
	ctx := c.Context
	conf := NewConfig(c)

	zlog.Info(ctx).Msg("🔥 Orchestrating the workload")
	err := orchestrateWorkload(ctx, conf)
	if err != nil {
		return err
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	zlog.Info(ctx).Stringer("duration", elapsedTime).Msg("Total time taken for completion")
	return nil
}

// orchestrateWorkload triggers the api endpoint hits and writes results to the desired location.
// It returns an error if any during the execution.
func orchestrateWorkload(ctx context.Context, conf *TestConfig) error {
	zlog.Info(ctx).Str("RunID", conf.RunID).Msg("Run details")
	var requests []map[string]interface{}
	var err error

	requests = attacker.CreateReadinessRequests(ctx, conf.HitSize, conf.Host)
	err = attacker.RunVegeta(ctx, requests, "get_readiness", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running GET operation on /readiness: %w", err)
	}

	requests = attacker.CreateLivenessRequests(ctx, conf.HitSize, conf.Host)
	err = attacker.RunVegeta(ctx, requests, "get_liveness", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running GET operation on /liveness: %w", err)
	}

	requests = attacker.CreateAuthorizedRequests(ctx, conf.HitSize, conf.Host, conf.AuthToken)
	err = attacker.RunVegeta(ctx, requests, "post_authorized", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running POST operation on /authorized: %w", err)
	}

	requests = attacker.CreateMetricsRequests(ctx, conf.HitSize, conf.Host, conf.AuthToken)
	err = attacker.RunVegeta(ctx, requests, "get_metrics", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running GET operation on /metrics: %w", err)
	}

	requests = attacker.CreateQueryRequests(ctx, conf.HitSize, conf.Host, conf.AuthToken, false)
	err = attacker.RunVegeta(ctx, requests, "post_query", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running POST operation on /v1/query: %w", err)
	}

	requests = attacker.CreateQueryRequests(ctx, conf.HitSize, conf.Host, conf.AuthToken, true)
	err = attacker.RunVegeta(ctx, requests, "post_query_with_cache", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running POST operation on /v1/query with cache: %w", err)
	}

	requests = attacker.CreateGetFeedbackStatusRequests(ctx, conf.HitSize, conf.Host, conf.AuthToken)
	err = attacker.RunVegeta(ctx, requests, "get_feedback_status", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running GET operation on /v1/feedback/status: %w", err)
	}

	requests = attacker.CreateFeedbackRequests(ctx, conf.HitSize, conf.Host, conf.AuthToken)
	err = attacker.RunVegeta(ctx, requests, "post_feedback", conf.Rps)
	if err != nil {
		return fmt.Errorf("Error while running POST operation on /v1/feedback: %w", err)
	}

	zlog.Info(ctx).Str("RunID", conf.RunID).Msg("👋 Exiting ols-load-generator")
	return nil
}
