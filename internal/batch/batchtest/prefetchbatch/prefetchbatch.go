package prefetchbatch

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/utils"
)

type prefetchExecuteRequest struct {
	id              string
	request         *PrefetchRequest
	mustWaitChan    *[]<-chan struct{}
	termChan        <-chan struct{}
	selfBroadCaster *utils.Broadcaster[struct{}]
}

func PrefetchBatch(ctr *app.Container, conf PrefetchConfig) (map[string]string, error) {
	store := sync.Map{}
	termBroadCaster := utils.NewBroadcaster[struct{}]()
	defer termBroadCaster.Close()

	requests := make([]prefetchExecuteRequest, len(conf.Requests))
	for i, req := range conf.Requests {
		requests[i] = prefetchExecuteRequest{
			id:              req.ID,
			request:         req,
			mustWaitChan:    &[]<-chan struct{}{},
			termChan:        termBroadCaster.Subscribe(),
			selfBroadCaster: utils.NewBroadcaster[struct{}](),
		}
		defer requests[i].selfBroadCaster.Close()
	}

	for _, req := range requests {
		if len(req.request.DependsOn) == 0 {
			continue
		}
		for _, depID := range req.request.DependsOn {
			var exists bool
			for _, depReq := range requests {
				if depID == depReq.id {
					*req.mustWaitChan = append(*req.mustWaitChan, depReq.selfBroadCaster.Subscribe())
					exists = true
				}
			}
			if !exists {
				return nil, fmt.Errorf("request %s depends on non-existent request(%s)", req.id, depID)
			}
		}
		ctr.Logger.Debug(ctr.Ctx, "request depends on",
			logger.Value("id", req.id), logger.Value("depends_on", req.request.DependsOn), logger.Value("on", "PrefetchBatch"))
	}

	atomicErr := atomic.Value{}
	wg := sync.WaitGroup{}
	for i, req := range requests {
		wg.Add(1)
		go func(preReq prefetchExecuteRequest) {
			defer func() {
				broadDone := preReq.selfBroadCaster.Broadcast(struct{}{})
				<-broadDone
				wg.Done()
			}()
			unprocessed := len(*preReq.mustWaitChan)
			for _, waitChan := range *preReq.mustWaitChan {
			loop:
				select {
				case <-preReq.termChan:
					<-waitChan
					return
				case <-waitChan:
					unprocessed--
					if unprocessed == 0 {
						break loop
					}
				}
			}
			err := executeRequest(ctr, i, preReq.request, preReq.termChan, &store, len(*preReq.mustWaitChan) > 0)

			if err != nil {
				atomicErr.Store(err)
				ctr.Logger.Error(ctr.Ctx, "failed to execute request",
					logger.Value("request_id", preReq.id), logger.Value("error", err), logger.Value("on", "PrefetchBatch"))
				ctr.Logger.Info(ctr.Ctx, "terminating all requests: failed to execute prefetch request",
					logger.Value("request_id", preReq.id), logger.Value("on", "PrefetchBatch"))
				termDoneChan := termBroadCaster.Broadcast(struct{}{})

				<-preReq.termChan
				<-termDoneChan
				return
			}
			ctr.Logger.Debug(ctr.Ctx, "request finished",
				logger.Value("id", preReq.id), logger.Value("on", "PrefetchBatch"))
		}(req)
	}

	ctr.Logger.Debug(ctr.Ctx, "waiting for all requests to finish",
		logger.Value("on", "PrefetchBatch"))
	wg.Wait()
	ctr.Logger.Debug(ctr.Ctx, "all requests finished",
		logger.Value("on", "PrefetchBatch"))

	if err := atomicErr.Load(); err != nil {
		ctr.Logger.Error(ctr.Ctx, "failed to find error",
			logger.Value("error", err.(error)), logger.Value("on", "PrefetchBatch"))
		return nil, err.(error)
	}

	result := make(map[string]string)
	store.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(string)
		return true
	})

	return result, nil
}
