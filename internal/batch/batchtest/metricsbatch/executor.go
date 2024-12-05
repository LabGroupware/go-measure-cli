package metricsbatch

import "os"

type MetricsThreadExecutor struct {
	outputFile *os.File
	fetcher    MetricsFetcher
}
