package batchtest

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type PipelineConfig struct {
	Type       string               `yaml:"type"`
	Concurrecy int                  `yaml:"concurrency"`
	Files      []PipelineFileConfig `yaml:"files"`
}

type PipelineFileConfig struct {
	ID   string `yaml:"id"`
	File string `yaml:"file"`
}

type executeRequest struct {
	id             string
	filename       string
	testRootDir    string
	metricsRootDir string
}

func pipelineBatch(ctx context.Context, ctr *app.Container, conf PipelineConfig, store *sync.Map, testOutput, metricsOutput string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	requests := make([]executeRequest, len(conf.Files))
	for i, f := range conf.Files {
		testDirPath := fmt.Sprintf("%s/%s", testOutput, f.ID)
		metricsDirPath := fmt.Sprintf("%s/%s", metricsOutput, f.ID)

		requests[i] = executeRequest{
			id:             f.ID,
			filename:       f.File,
			testRootDir:    testDirPath,
			metricsRootDir: metricsDirPath,
		}
	}

	var sequential bool
	concurrency := conf.Concurrecy
	if concurrency < 0 {
		concurrency = len(requests)
	}
	if concurrency == 0 {
		concurrency = 1
		sequential = true
	}

	if sequential {
		for _, req := range requests {

			err := baseExecute(ctx, ctr, req.filename, store, req.testRootDir, req.metricsRootDir)

			if err != nil {
				ctr.Logger.Error(ctr.Ctx, "failed to execute request",
					logger.Value("error", err), logger.Value("on", "PipelineBatch"))
				return fmt.Errorf("failed to execute request: %v", err)
			}

			ctr.Logger.Debug(ctr.Ctx, "request finished",
				logger.Value("on", "PipelineBatch"))
		}

		return nil
	} else {
		atomicErr := atomic.Value{}
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrency)
		for _, req := range requests {
			wg.Add(1)

			go func(preReq executeRequest) {
				defer wg.Done()

				sem <- struct{}{}

				err := baseExecute(ctx, ctr, preReq.filename, store, preReq.testRootDir, preReq.metricsRootDir)

				if err != nil {
					atomicErr.Store(err)
					ctr.Logger.Error(ctr.Ctx, "failed to execute request",
						logger.Value("error", err), logger.Value("on", "PipelineBatch"))
					cancel()
					return
				}
				ctr.Logger.Debug(ctr.Ctx, "request finished",
					logger.Value("on", "PipelineBatch"))

				<-sem
			}(req)
		}

		wg.Wait()

		close(sem)

		if err := atomicErr.Load(); err != nil {
			ctr.Logger.Error(ctr.Ctx, "failed to find error",
				logger.Value("error", err.(error)), logger.Value("on", "PipelineBatch"))
			return err.(error)
		}

		return nil
	}
}
