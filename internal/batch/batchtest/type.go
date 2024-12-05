package batchtest

import (
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/metricsbatch"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/prefetchbatch"
)

type BatchTestType struct {
	Type     string                          `yaml:"type"`
	Prefetch prefetchbatch.PrefetchConfig    `yaml:"prefetch"`
	Metrics  metricsbatch.MetricsBatchConfig `yaml:"metrics"`
	Data     any                             `yaml:"data"`
	Output   BatchTestOutput                 `yaml:"output"`
}

type BatchTestOutput struct {
	Enabled bool `yaml:"enabled"`
}
