package massquery

import "sync"

type MassiveQueryExecutor struct {
	threadExecutors []*MassiveQueryThreadExecutor
}

func NewMassiveQueryExecutor() *MassiveQueryExecutor {
	return &MassiveQueryExecutor{}
}

func (e *MassiveQueryExecutor) Execute() {

	wg := sync.WaitGroup{}

	for _, threadExecutor := range e.threadExecutors {
		wg.Add(1)
		go func() {
			defer wg.Done()
			threadExecutor.Execute()
		}()
	}

	wg.Wait()
}
