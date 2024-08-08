package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-bulldozer/go-commons/indexers"
	ocpmetadata "github.com/cloud-bulldozer/go-commons/ocp-metadata"
	"github.com/google/uuid"
	"github.com/kube-burner/kube-burner/pkg/burner"
	"github.com/kube-burner/kube-burner/pkg/config"
	"github.com/kube-burner/kube-burner/pkg/prometheus"
	"github.com/kube-burner/kube-burner/pkg/util/metrics"
	"github.com/kube-burner/kube-burner/pkg/workloads"
	"github.com/quay/zlog"
	"github.com/urfave/cli/v2"
)

// IndexCmd handles createtoken CLI.
var IndexCmd = &cli.Command{
	Name:        "index",
	Description: "Indexes metrics within given timerange",
	Usage:       "ols-load-generator index",
	Action:      IndexAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "uuid",
			Usage:   "--uuid f519d9b2-aa62-44ab-9ce8-4156b712f6d2",
			Value:   uuid.New().String(),
			EnvVars: []string{"OLS_TEST_UUID"},
		},
		&cli.StringFlag{
			Name:    "eshost",
			Usage:   "--eshost eshosturl",
			Value:   "",
			EnvVars: []string{"OLS_TEST_ES_HOST"},
		},
		&cli.StringFlag{
			Name:    "esindex",
			Usage:   "--esindex esindex",
			Value:   "",
			EnvVars: []string{"OLS_TEST_ES_INDEX"},
		},
		&cli.IntFlag{
			Name:    "metricstep",
			Usage:   "--metricstep 30",
			Value:   30,
			EnvVars: []string{"OLS_TEST_METRIC_STEP"},
		},
		&cli.Int64Flag{
			Name:    "start",
			Usage:   "--start 1720410990",
			Value:   time.Now().Unix() - 3600,
			EnvVars: []string{"OLS_TEST_START"},
		},
		&cli.Int64Flag{
			Name:    "end",
			Usage:   "--end 1720470990",
			Value:   time.Now().Unix(),
			EnvVars: []string{"OLS_TEST_END"},
		},
		&cli.StringFlag{
			Name:    "profiles",
			Usage:   "--profiles metrics.yaml,metrics-report.yaml",
			Value:   "attacker/assets/profiles/metrics-report.yaml,attacker/assets/profiles/metrics-timeseries.yaml",
			EnvVars: []string{"OLS_TEST_PROFILES"},
		},
	},
}

// indexConfig creates and returns a test configuration from CLI options.
func indexConfig(c *cli.Context) *TestConfig {
	profilesArg := c.String("profiles")
	return &TestConfig{
		Uuid:       c.String("uuid"),
		ESHost:     c.String("eshost"),
		ESIndex:    c.String("esindex"),
		MetricStep: c.Int("metricstep"),
		Profiles:   strings.Split(strings.TrimSpace(profilesArg), ","),
	}
}

// IndexAction allows us to trigger index action.
// It returns an error if any during the execution.
func IndexAction(c *cli.Context) error {
	startTime := time.Now()
	ctx := c.Context
	conf := indexConfig(c)

	zlog.Info(ctx).Msg("📁 Indexing metrics")
	err := IndexMetrics(ctx, conf, time.Unix(c.Int64("start"), 0).Add(-10*time.Minute), time.Unix(c.Int64("end"), 0).Add(10*time.Minute))
	if err != nil {
		return err
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	zlog.Info(ctx).Stringer("duration", elapsedTime).Msg("Total time taken for completion")
	return nil
}

// IndexMetrics allow us to index prometheus metrics.
// It returns an error if any during the execution.
func IndexMetrics(ctx context.Context, conf *TestConfig, startTime, endTime time.Time) error {
	var err error
	var prometheusURL, prometheusToken string
	var metadataAgent ocpmetadata.Metadata
	var clusterMetadata ocpmetadata.ClusterMetadata
	metadata := make(map[string]interface{})
	kubeClientProvider := config.NewKubeClientProvider("", "")
	_, restConfig := kubeClientProvider.DefaultClientSet()
	metadataAgent, err = ocpmetadata.NewMetadata(restConfig)
	if err != nil {
		return err
	}
	workloads.ConfigSpec.GlobalConfig.UUID = conf.Uuid
	prometheusURL, prometheusToken, err = metadataAgent.GetPrometheus()
	if err != nil {
		return fmt.Errorf("Error obtaining prometheus information: %v", err)
	}
	clusterMetadata, err = metadataAgent.GetClusterMetadata()
	if err != nil {
		return err
	}
	jsonData, _ := json.Marshal(clusterMetadata)
	json.Unmarshal(jsonData, &metadata)
	indexer := config.MetricsEndpoint{
		Endpoint:      prometheusURL,
		Token:         prometheusToken,
		Step:          time.Duration(conf.MetricStep) * time.Second,
		Metrics:       conf.Profiles,
		SkipTLSVerify: true,
	}
	if conf.ESHost != "" && conf.ESIndex != "" {
		indexer.IndexerConfig = indexers.IndexerConfig{
			Type:               indexers.OpenSearchIndexer,
			Servers:            []string{conf.ESHost},
			Index:              conf.ESIndex,
			InsecureSkipVerify: true,
		}
	} else {
		indexer.IndexerConfig = indexers.IndexerConfig{
			Type:             indexers.LocalIndexer,
			MetricsDirectory: "collected-metrics" + "-" + conf.Uuid,
		}
	}

	workloads.ConfigSpec.MetricsEndpoints = append(workloads.ConfigSpec.MetricsEndpoints, indexer)
	metricsScraper := metrics.ProcessMetricsScraperConfig(metrics.ScraperConfig{
		ConfigSpec: &workloads.ConfigSpec,
	})
	metricsScraper.Metadata = nil
	for _, prometheusClient := range metricsScraper.PrometheusClients {
		prometheusJob := prometheus.Job{
			Start: time.Unix(startTime.Unix(), int64(startTime.Nanosecond())),
			End:   time.Unix(endTime.Unix(), int64(endTime.Nanosecond())),
			JobConfig: config.Job{
				Name: "ols-load-generator",
			},
		}
		if prometheusClient.ScrapeJobsMetrics(prometheusJob) != nil {
			return fmt.Errorf("Error while scraping cluster prometheus")
		}
	}
	var indexerValue indexers.Indexer
	for _, value := range metricsScraper.IndexerList {
		indexerValue = value
		break
	}
	jobSummary := burner.JobSummary{
		Timestamp:    startTime.UTC(),
		EndTimestamp: endTime.UTC(),
		ElapsedTime:  endTime.UTC().Sub(startTime.UTC()).Round(time.Second).Seconds(),
		UUID:         conf.Uuid,
		JobConfig: config.Job{
			Name: "ols-load-generator",
		},
		Metadata:   metadata,
		MetricName: "jobSummary",
		Passed:     true,
	}
	burner.IndexJobSummary([]burner.JobSummary{jobSummary}, indexerValue)
	return nil
}
