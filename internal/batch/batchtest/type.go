package batchtest

import "github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/prefetchbatch"

type BatchTestType struct {
	Type     string                       `yaml:"type"`
	Prefetch prefetchbatch.PrefetchConfig `yaml:"prefetch"`
	Data     any                          `yaml:"data"`
}
