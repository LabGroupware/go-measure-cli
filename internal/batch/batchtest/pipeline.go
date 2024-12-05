package batchtest

import (
	"context"
	"fmt"
	"os"
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
	var err error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("Pipeline called!!", testOutput, metricsOutput)

	requests := make([]executeRequest, len(conf.Files))
	for i, f := range conf.Files {
		testDirPath := fmt.Sprintf("%s/%s", testOutput, f.ID)
		err = os.MkdirAll(testDirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		metricsDirPath := fmt.Sprintf("%s/%s", metricsOutput, f.ID)
		err = os.MkdirAll(metricsDirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

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

			fmt.Println("Request finished on seq", testOutput, metricsOutput, req.id, err)

			if err != nil {
				ctr.Logger.Error(ctr.Ctx, "failed to execute request",
					logger.Value("error", err), logger.Value("on", "PipelineBatch"))
				return fmt.Errorf("failed to execute request: %v", err)
			}

			ctr.Logger.Debug(ctr.Ctx, "request finished",
				logger.Value("on", "PipelineBatch"))
		}

		fmt.Println("All requests finished on seq", testOutput, metricsOutput)

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

				fmt.Println("Request finished on parallel", testOutput, metricsOutput, req.id, err)

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

		fmt.Println("All requests finished on parallel", testOutput, metricsOutput, atomicErr.Load())

		close(sem)

		if err := atomicErr.Load(); err != nil {
			ctr.Logger.Error(ctr.Ctx, "failed to find error",
				logger.Value("error", err.(error)), logger.Value("on", "PipelineBatch"))
			return err.(error)
		}

		return nil
	}
}
