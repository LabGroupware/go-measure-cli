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

type FindFileObjectFactory struct{}

func (f FindFileObjectFactory) Factory(
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	authToken *auth.AuthToken,
	apiEndpoint string,
	outputFile *os.File,
) (queryreq.QueryExecutor, error) {
	var ok bool
	var fileObjectId string

	req := queryreq.FindFileObjectReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	if fileObjectId, ok = request.PathVariables["fileObjectId"]; !ok {
		return nil, fmt.Errorf("fileObjectId not found in pathVariables")
	}
	if fileObjectId == "*" {
		fileObjectId = testprompt.GenerateRandomString(10)
	}
	req.Path.FileObjectID = fileObjectId

	for key, param := range request.QueryParam {
		switch key {
		case "with":
			req.Param.With = param
		}
	}

	resChan := make(chan queryreq.ResponseContent[response.ResponseDto[domain.FileObjectDto]])

	runAsyncProcessing(ctr, id, request, termChan, resChan, outputFile)

	return queryreq.RequestContent[queryreq.FindFileObjectReq, response.ResponseDto[domain.FileObjectDto]]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
	}, nil
}

type GetFileObjectsFactory struct{}

func (f GetFileObjectsFactory) Factory(
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	authToken *auth.AuthToken,
	apiEndpoint string,
	outputFile *os.File,
) (queryreq.QueryExecutor, error) {
	req := queryreq.GetFileObjectsReq{
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
		case "hasBucketFilter":
			hasBucketFilterBool, err := strconv.ParseBool(param[0])
			if err != nil {
				ctr.Logger.Warn(ctr.Ctx, "Failed to convert hasBucketFilter to bool",
					logger.Value("error", err))
				continue
			}
			req.Param.HasBucketFilter = hasBucketFilterBool
		case "filterBucketIDs":
			req.Param.FilterBucketIDs = param
		case "with":
			req.Param.With = param
		}

	}

	resChan := make(chan queryreq.ResponseContent[response.ListResponseDto[domain.FileObjectDto]])

	runAsyncProcessing(ctr, id, request, termChan, resChan, outputFile)

	return queryreq.RequestContent[queryreq.GetFileObjectsReq, response.ListResponseDto[domain.FileObjectDto]]{
		Req:          req,
		Interval:     request.Interval,
		ResponseWait: request.AwaitPrevResp,
		ResChan:      resChan,
	}, nil
}
