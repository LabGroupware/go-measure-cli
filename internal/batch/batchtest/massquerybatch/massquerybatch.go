package massquerybatch

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/queryreqbatch"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

func MassQueryBatch(ctr *app.Container, massQuery MassQuery) error {
	var err error
	concurrentCount := len(massQuery.Data.Requests)

	threadExecutors := make([]*MassiveQueryThreadExecutor, concurrentCount)

	timestamp := time.Now().Format("20060102_150405")
	dirPath := fmt.Sprintf("%s/test_%s", ctr.Config.Batch.Test.MassQuery.Output, timestamp)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	for i := 0; i < concurrentCount; i++ {
		request := massQuery.Data.Requests[i]
		logFilePath := fmt.Sprintf("%s/logcsv_%010d.csv", dirPath, i+1)
		file, err := os.Create(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		threadExecutors[i] = &MassiveQueryThreadExecutor{
			ID:         i + 1,
			outputFile: file,
		}
		defer threadExecutors[i].Close(ctr.Ctx)

		writer := csv.NewWriter(file)
		header := []string{"Success", "SendDatetime", "ReceivedDatetime", "Count", "ResponseTime", "StatusCode", "Data"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}
		writer.Flush()

		endType := massQuery.Data.Requests[i].EndpointType
		queryType := queryreqbatch.NewQueryTypeFromString(endType)
		factor := queryreqbatch.TypeFactoryMap[queryType]
		termChan := make(chan queryreqbatch.TerminateType)
		defer close(termChan)
		validatedReq := &queryreqbatch.ValidatedQueryRequest{}
		if err := queryreqbatch.ValidateQueryReq(ctr, request, validatedReq); err != nil {
			return fmt.Errorf("failed to validate query request: %v", err)
		}
		writeFunc := func(
			ctr *app.Container,
			id int,
			data queryreqbatch.WriteData,
		) error {
			writer := csv.NewWriter(threadExecutors[i].outputFile)
			ctr.Logger.Debug(ctr.Ctx, "Writing data to csv",
				logger.Value("id", id), logger.Value("data", data), logger.Value("on", "runAsyncProcessing"))
			if err := writer.Write(data.ToSlice()); err != nil {
				ctr.Logger.Error(ctr.Ctx, "failed to write data to csv",
					logger.Value("error", err), logger.Value("on", "runAsyncProcessing"))
			}
			writer.Flush()
			return nil
		}
		executor, resCloser, err := factor.Factory(
			ctr,
			i+1,
			validatedReq,
			termChan,
			ctr.AuthToken,
			ctr.Config.Web.API.Url,
			writeFunc,
		)
		ctr.Logger.Info(ctr.Ctx, "created executor",
			logger.Value("id", i+1), logger.Value("type", endType), logger.Value("executor", executor))
		if err != nil {
			return fmt.Errorf("failed to create executor: %v", err)
		}
		threadExecutors[i].RequestExecutor = executor
		threadExecutors[i].TermChan = termChan
		threadExecutors[i].responseChanCloser = resCloser
	}

	var wg sync.WaitGroup
	var startChan = make(chan struct{})
	for _, executor := range threadExecutors {
		wg.Add(1)
		go func(exec *MassiveQueryThreadExecutor) {
			defer wg.Done()
			defer exec.responseChanCloser()
			if err := exec.Execute(ctr, startChan); err != nil {
				ctr.Logger.Error(ctr.Ctx, "failed to execute query",
					logger.Value("error", err), logger.Value("id", exec.ID))
			}
		}(executor)
	}
	close(startChan)
	wg.Wait()

	return nil
}
