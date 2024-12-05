package queryreqbatch

import (
	"context"
	"fmt"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/api/response"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
)

type FindJobFactory struct{}

func (f FindJobFactory) Factory(
	ctx context.Context,
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- TerminateType,
	authToken *auth.AuthToken,
	apiEndpoint string,
	consumer ResponseDataConsumer,
) (queryreq.QueryExecutor, func(), error) {
	var ok bool
	var jobId string

	req := queryreq.GetJobReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	if jobId, ok = request.PathVariables["jobId"]; !ok {
		return nil, nil, fmt.Errorf("jobId not found in pathVariables")
	}
	if jobId == "*" {
		jobId = testprompt.GenerateRandomString(10)
	}
	req.Path.JobID = jobId

	// INFO: close on executor, because only it will write to this channel
	resChan := make(chan queryreq.ResponseContent[response.JobResponseDto])

	resChanCloser := func() {
		close(resChan)
	}

	runAsyncProcessing(ctx, ctr, id, request, termChan, resChan, consumer)

	return queryreq.RequestContent[queryreq.GetJobReq, response.JobResponseDto]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
		CountLimit:   request.Break.Count,
	}, resChanCloser, nil
}
