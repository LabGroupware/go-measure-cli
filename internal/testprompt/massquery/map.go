package massquery

import (
	"context"
	"os"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/executor"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
)

type ExecutorFactory interface {
	factory(
		ctx context.Context,
		id int,
		interval time.Duration,
		responseWait bool,
		termChan chan<- struct{},
		authToken *auth.AuthToken,
		apiEndpoint string,
		outputFile *os.File,
	) executor.RequestExecutor
}

var typeFactoryMap = map[QueryType]ExecutorFactory{
	FindJob: FindJobFactory{},
}
