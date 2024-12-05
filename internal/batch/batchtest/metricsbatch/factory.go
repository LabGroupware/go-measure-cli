package metricsbatch

import (
	"context"

	"github.com/LabGroupware/go-measure-tui/internal/app"
)

type MetricsFetcherFactory interface {
	FetcherFactory(ctx context.Context, ctr *app.Container) (MetricsFetcher, error)
}

var metricsFetcherFactoryMap = map[MetricsType]MetricsFetcherFactory{
	MetricsTypePrometheus: &PrometheusMetricsBatchRequestConfig{},
}
