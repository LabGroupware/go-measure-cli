package metricsbatch

import (
	"context"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
)

type MetricsBatchConfig struct {
	Enabled  bool                        `yaml:"enabled"`
	Requests []MetricsBatchRequestConfig `yaml:"requests"`
}

type MetricsBatchRequestConfig struct {
	ID       string                          `yaml:"id"`
	Type     string                          `yaml:"type"`
	Data     []MetricsBatchRequestDataConfig `yaml:"data"`
	Interval string                          `yaml:"interval"`
}

func (m *MetricsBatchRequestConfig) Validate(ctx context.Context, ctr *app.Container, validated *ValidatedMetricsBatchRequestConfig) error {
	validated.ID = m.ID
	return nil
}

type ValidatedMetricsBatchRequestConfig struct {
	ID       string
	Type     MetricsType
	Data     []MetricsBatchRequestDataConfig
	Interval time.Duration
}

type MetricsBatchRequestDataConfig struct {
	Key      string `yaml:"key"`
	JmesPath string `yaml:"jmesPath"`
}

type MetricsType string

const (
	MetricsTypePrometheus MetricsType = "prometheus"
)

type MetricsFetcher interface {
	Fetch(ctx context.Context, ctr *app.Container) (any, chan<- struct{}, error)
}

type TermType int

const (
	_ TermType = iota
	TermTypeContext
	TermTypeTerm
)
