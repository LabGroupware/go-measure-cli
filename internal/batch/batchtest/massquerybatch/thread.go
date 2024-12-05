package massquerybatch

import (
	"context"
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

	<-startChan

	ctr.Logger.Info(ctx, "Query Start",
		logger.Value("QueryID", e.ID), logger.Value("OutputFile", e.outputFile.Name()))
	term, err := e.RequestExecutor.QueryExecute(ctx, ctr)
	if err != nil {
		return err
	}

	<-e.TermChan
	term <- struct{}{}
	ctr.Logger.Info(ctx, "Query End For Term",
		logger.Value("QueryID", e.ID), logger.Value("OutputFile", e.outputFile.Name()))
	return nil
}

func (e *MassiveQueryThreadExecutor) Close(ctx context.Context) error {
	e.outputFile.Close()
	return nil
}
