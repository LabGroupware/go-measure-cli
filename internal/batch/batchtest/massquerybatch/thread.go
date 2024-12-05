package massquerybatch

import (
	"context"
	"fmt"
	"os"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/queryreqbatch"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type MassiveQueryThreadExecutor struct {
	ID              int
	outputFile      *os.File
	RequestExecutor queryreq.QueryExecutor
	TermChan        chan queryreqbatch.TerminateType
	successBreak    []string
}

func NewMassiveQueryThreadExecutor(
	id int,
	outputFile *os.File,
	req queryreq.QueryExecutor,
) *MassiveQueryThreadExecutor {
	return &MassiveQueryThreadExecutor{
		ID:              id,
		outputFile:      outputFile,
		RequestExecutor: req,
	}
}

func (e *MassiveQueryThreadExecutor) Execute(
	ctx context.Context,
	ctr *app.Container,
	startChan <-chan struct{},
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	<-startChan

	ctr.Logger.Info(ctx, "Query Start",
		logger.Value("QueryID", e.ID))
	err := e.RequestExecutor.QueryExecute(ctx, ctr)
	if err != nil {
		return err
	}

	termType := <-e.TermChan
	ctr.Logger.Info(ctx, "Query End For Break",
		logger.Value("QueryID", e.ID))
	for _, breakType := range e.successBreak {
		if termType == queryreqbatch.NewTerminateTypeFromString(breakType) {
			ctr.Logger.Info(ctx, "Query End For Success Break", logger.Value("QueryID", e.ID))
			return nil
		}
	}
	return fmt.Errorf("query End For Fail Break: %s", termType.String())
}

func (e *MassiveQueryThreadExecutor) Close(ctx context.Context) error {
	e.outputFile.Close()
	return nil
}
