package queryreqbatch

import (
	"fmt"
	"os"
	"strconv"

	"github.com/LabGroupware/go-measure-tui/internal/api/domain"
	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/api/response"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
)

type FindUserFactory struct{}

func (f FindUserFactory) Factory(
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	authToken *auth.AuthToken,
	apiEndpoint string,
	outputFile *os.File,
) (queryreq.QueryExecutor, error) {
	var ok bool
	var userId string

	req := queryreq.FindUserReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	if userId, ok = request.PathVariables["userId"]; !ok {
		return nil, fmt.Errorf("userId not found in pathVariables")
	}
	if userId == "*" {
		userId = testprompt.GenerateRandomString(10)
	}
	req.Path.UserID = userId

	for key, param := range request.QueryParam {
		switch key {
		case "with":
			req.Param.With = param
		}
	}

	resChan := make(chan queryreq.ResponseContent[response.ResponseDto[domain.UserProfileDto]])

	runAsyncProcessing(ctr, id, request, termChan, resChan, outputFile)

	return queryreq.RequestContent[queryreq.FindUserReq, response.ResponseDto[domain.UserProfileDto]]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
	}, nil
}

type GetUsersFactory struct{}

func (f GetUsersFactory) Factory(
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	authToken *auth.AuthToken,
	apiEndpoint string,
	outputFile *os.File,
) (queryreq.QueryExecutor, error) {
	req := queryreq.GetUsersReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	for key, param := range request.QueryParam {
		switch key {
		case "limit":
			limitInt, err := strconv.Atoi(param[0])
			if err != nil {
				ctr.Logger.Warn(ctr.Ctx, "Failed to convert limit to int",
					logger.Value("error", err))
				continue
			}
			req.Param.Limit = limitInt
		case "offset":
			offsetInt, err := strconv.Atoi(param[0])
			if err != nil {
				ctr.Logger.Warn(ctr.Ctx, "Failed to convert offset to int",
					logger.Value("error", err))
				continue
			}
			req.Param.Offset = offsetInt
		case "cursor":
			req.Param.Cursor = param[0]
		case "pagination":
			req.Param.Pagination = param[0]
		case "sortField":
			req.Param.SortField = param[0]
		case "sortOrder":
			req.Param.SortOrder = param[0]
		case "withCount":
			withCountBool, err := strconv.ParseBool(param[0])
			if err != nil {
				ctr.Logger.Warn(ctr.Ctx, "Failed to convert withCount to bool",
					logger.Value("error", err))
				continue
			}
			req.Param.WithCount = withCountBool
		case "with":
			req.Param.With = param
		}

	}

	resChan := make(chan queryreq.ResponseContent[response.ListResponseDto[domain.UserProfileDto]])

	runAsyncProcessing(ctr, id, request, termChan, resChan, outputFile)

	return queryreq.RequestContent[queryreq.GetUsersReq, response.ListResponseDto[domain.UserProfileDto]]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
	}, nil
}
