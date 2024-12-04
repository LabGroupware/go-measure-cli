package queryreqbatch

import (
	"fmt"

	"github.com/LabGroupware/go-measure-tui/internal/api/domain"
	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/api/response"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
)

type FindUserPreferenceFactory struct{}

func (f FindUserPreferenceFactory) Factory(
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- TerminateType,
	authToken *auth.AuthToken,
	apiEndpoint string,
	consumer ResponseDataConsumer,
) (queryreq.QueryExecutor, func(), error) {
	var ok bool
	var userPreferenceId string

	req := queryreq.FindUserPreferenceReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	if userPreferenceId, ok = request.PathVariables["userPreferenceId"]; !ok {
		return nil, nil, fmt.Errorf("userPreferenceId not found in pathVariables")
	}
	if userPreferenceId == "*" {
		userPreferenceId = testprompt.GenerateRandomString(10)
	}
	req.Path.UserPreferenceID = userPreferenceId

	for key, param := range request.QueryParam {
		switch key {
		case "with":
			req.Param.With = param
		}
	}

	resChan := make(chan queryreq.ResponseContent[response.ResponseDto[domain.UserPreferenceDto]])

	resChanCloser := func() {
		close(resChan)
	}

	runAsyncProcessing(ctr, id, request, termChan, resChan, consumer)

	return queryreq.RequestContent[queryreq.FindUserPreferenceReq, response.ResponseDto[domain.UserPreferenceDto]]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
		CountLimit:   request.Break.Count,
	}, resChanCloser, nil
}
