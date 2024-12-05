package massquery

import (
	"context"
	"os"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/executor"
)

type MassiveQueryThreadExecutor struct {
	ID              int
	outputFile      *os.File
	RequestExecutor executor.RequestExecutor
	TermChan        chan struct{}
}

func NewMassiveQueryThreadExecutor(
	id int,
	outputFile *os.File,
	req executor.RequestExecutor,
) *MassiveQueryThreadExecutor {
	return &MassiveQueryThreadExecutor{
		ID:              id,
		outputFile:      outputFile,
		RequestExecutor: req,
	}
}

func (e *MassiveQueryThreadExecutor) Execute(ctx context.Context) error {
	// fmt.Printf("Query %d の実行を開始します\n", e.ID)
	// term, err := e.RequestExecutor.QueryExecute(ctx)
	// if err != nil {
	// 	return err
	// }

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return nil
	// 	case <-e.TermChan:
	// 		term <- struct{}{}
	// 		fmt.Printf("終了条件を満たしたため, Query %d の実行を終了します\n", e.ID)
	// 		return nil
	// 	}
	// }
	return nil
}

func (e *MassiveQueryThreadExecutor) Close(ctx context.Context) error {
	e.outputFile.Close()

	return nil
}
