package queryreqbatch

import (
	"fmt"
	"os"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/api/response"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
)

type FindJobFactory struct{}

func (f FindJobFactory) Factory(
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	authToken *auth.AuthToken,
	apiEndpoint string,
	outputFile *os.File,
) (queryreq.QueryExecutor, error) {
	var ok bool
	var jobId string

	req := queryreq.GetJobReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	if jobId, ok = request.PathVariables["jobId"]; !ok {
		return nil, fmt.Errorf("jobId not found in pathVariables")
	}
	if jobId == "*" {
		jobId = testprompt.GenerateRandomString(10)
	}
	req.Path.JobID = jobId

	resChan := make(chan queryreq.ResponseContent[response.JobResponseDto])

	runAsyncProcessing(ctr, id, request, termChan, resChan, outputFile)

	return queryreq.RequestContent[queryreq.GetJobReq, response.JobResponseDto]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
	}, nil
}
