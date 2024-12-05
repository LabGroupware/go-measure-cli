package waitsagabatch

import (
	"context"
	"fmt"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

func WaitSagaBatch(
	ctx context.Context,
	ctr *app.Container,
	// conf OneExecuteConfig,
	store *sync.Map,
) (map[string]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sock, err := NewSocket(ctx, ctr)
	if err != nil {
		ctr.Logger.Error(ctx, "failed to create socket",
			logger.Value("error", err))
		return nil, fmt.Errorf("failed to create socket: %v", err)
	}

	done, err := sock.Connect(ctx, ctr)
	if err != nil {
		ctr.Logger.Error(ctx, "failed to connect to socket",
			logger.Value("error", err))
		return nil, fmt.Errorf("failed to connect to socket: %v", err)
	}

	select {
	case <-done:
		ctr.Logger.Warn(ctx, "socket connection closed")
		return nil, fmt.Errorf("socket connection closed")
	case <-ctx.Done():
		ctr.Logger.Warn(ctx, "context cancelled")
		return nil, fmt.Errorf("context cancelled")
	}

	// newStore := sync.Map{}
	// err := executeRequest(ctx, ctr, 0, conf.Request, &newStore)
	// if err != nil {
	// 	ctr.Logger.Error(ctx, "failed to find error",
	// 		logger.Value("error", err))
	// 	return nil, fmt.Errorf("failed to execute request: %v", err)
	// }

	// newMap := make(map[string]string)

	// newStore.Range(func(key, value interface{}) bool {
	// 	store.Store(key, value)
	// 	newMap[key.(string)] = value.(string)
	// 	return true
	// })

	// return newMap, nil
	return nil, nil
}
