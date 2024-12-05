package batchtest

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/massquerybatch"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/metricsbatch"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/onequerybatch"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/prefetchbatch"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"gopkg.in/yaml.v3"
)

func baseExecute(
	ctx context.Context,
	ctr *app.Container,
	filename string,
	store *sync.Map,
	outputRoot string,
	metricsOutputRoot string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	file, err := os.Open(filepath.Join(ctr.Config.Batch.Test.Path, filename))
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var conf BatchTestType
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		return fmt.Errorf("failed to decode yaml: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	if conf.Prefetch.Enabled {
		var replacements = make(map[string]string)
		if replacements, err = prefetchbatch.PrefetchBatch(ctx, ctr, conf.Prefetch, store); err != nil {
			return fmt.Errorf("failed to execute prefetch: %v", err)
		}

		ctr.Logger.Debug(ctx, "replacements set",
			logger.Value("replacements", replacements))
	}

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {

		return fmt.Errorf("failed to read file: %w", err)
	}

	content := buffer.String()
	placeholderRegex := regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

	result := placeholderRegex.ReplaceAllStringFunc(content, func(match string) string {
		key := placeholderRegex.FindStringSubmatch(match)[1]
		if v, exists := store.Load(key); exists {
			return v.(string)
		}
		return match
	})

	var yamlData map[string]interface{}

	if err := yaml.Unmarshal([]byte(result), &yamlData); err != nil {
		return fmt.Errorf("failed to parse as YAML: %w", err)
	}

	ctr.Logger.Debug(ctx, "replaced content",
		logger.Value("content", yamlData))

	reader := bytes.NewReader([]byte(result))

	if conf.Metrics.Enabled {
		ctr.Logger.Debug(ctx, "metrics enabled",
			logger.Value("metrics", conf.Metrics))
		err = os.MkdirAll(metricsOutputRoot, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		if err := metricsbatch.MetricsFetchBatch(
			ctx,
			ctr,
			conf.Metrics,
			bytes.NewReader([]byte(result)),
			conf.Type,
			metricsOutputRoot,
		); err != nil {
			return fmt.Errorf("failed to execute metrics fetch: %v", err)
		}
	}

	switch conf.Type {
	case "MassQuery", "Pipeline":
		err := os.MkdirAll(outputRoot, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	switch conf.Type {
	case "OneQuery":
		var oneQuery onequerybatch.OneQueryConfig
		decoder := yaml.NewDecoder(reader)
		if err := decoder.Decode(&oneQuery); err != nil {
			return fmt.Errorf("failed to decode yaml: %v", err)
		}
		var values map[string]string
		if values, err = onequerybatch.OneQueryBatch(ctx, ctr, oneQuery, store); err != nil {
			return fmt.Errorf("failed to execute one query: %v", err)
		}
		store.Range(func(key, value interface{}) bool {
			ctr.Logger.Debug(ctx, "current store value",
				logger.Value("key", key), logger.Value("value", value))
			return true
		})
		ctr.Logger.Info(ctx, "newValues",
			logger.Value("values", values))
	case "MassQuery":
		var massQuery massquerybatch.MassQuery
		decoder := yaml.NewDecoder(reader)
		if err := decoder.Decode(&massQuery); err != nil {
			return fmt.Errorf("failed to decode yaml: %v", err)
		}
		if err := massquerybatch.MassQueryBatch(ctx, ctr, massQuery, outputRoot); err != nil {
			return fmt.Errorf("failed to execute mass query: %v", err)
		}
	case "WaitSaga":

		return fmt.Errorf("not implemented")
	case "Pipeline":
		var pipeline PipelineConfig
		decoder := yaml.NewDecoder(reader)
		if err := decoder.Decode(&pipeline); err != nil {

			return fmt.Errorf("failed to decode yaml: %v", err)
		}
		if err := pipelineBatch(ctx, ctr, pipeline, store, outputRoot, metricsOutputRoot); err != nil {
			return fmt.Errorf("failed to execute pipeline: %v", err)
		}
	default:

		return fmt.Errorf("unknown type")
	}

	return nil
}
