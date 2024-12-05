package metricsbatch

import (
	"context"

	"github.com/LabGroupware/go-measure-tui/internal/app"
)

type MetricsFetcherFactory interface {
	Init(conf []byte) error
	FetcherFactory(ctx context.Context, ctr *app.Container) (MetricsFetcher, error)
}

var metricsFetcherFactoryMap = map[MetricsType]MetricsFetcherFactory{
	MetricsTypePrometheus: &PrometheusMetricsBatchRequestConfig{},
}
