package massquery

import (
	"context"
	"sync"
)

type MassiveQueryExecutor struct {
	threadExecutors []*MassiveQueryThreadExecutor
}

func NewMassiveQueryExecutor() *MassiveQueryExecutor {
	return &MassiveQueryExecutor{
		threadExecutors: []*MassiveQueryThreadExecutor{},
	}
}

func NewMassiveQueryExecutorWithThreads(threads []*MassiveQueryThreadExecutor) *MassiveQueryExecutor {
	return &MassiveQueryExecutor{
		threadExecutors: threads,
	}
}

func (e *MassiveQueryExecutor) Execute(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for _, threadExecutor := range e.threadExecutors {
		wg.Add(1)
		go func() {
			defer wg.Done()
			threadExecutor.Execute(ctx)
		}()
	}

	wg.Wait()

	return nil
}

func (e *MassiveQueryExecutor) Close(ctx context.Context) {
	for _, threadExecutor := range e.threadExecutors {
		threadExecutor.Close(ctx)
	}
}
