package metricsbatch

import (
	"context"

	"github.com/LabGroupware/go-measure-tui/internal/app"
)

func PrefetchBatch(ctx context.Context, ctr *app.Container, conf MetricsBatchConfig, testType string, outputRoot string) (map[string]string, error) {

	// var err error
	// concurrentCount := len(conf.Requests)

	// threadExecutors := make([]*MetricsThreadExecutor, concurrentCount)

	// timestamp := time.Now().Format("20060102_150405")
	// dirPath := fmt.Sprintf("%s/%s/test_%s", ctr.Config.Batch.Metrics.Output, testType, timestamp)
	// err = os.MkdirAll(dirPath, os.ModePerm)
	// if err != nil {
	// 	return fmt.Errorf("failed to create directory: %v", err)
	// }
	return nil, nil
}
