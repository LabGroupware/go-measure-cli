package prefetchbatch

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/queryreqbatch"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/utils"
	"github.com/jmespath/go-jmespath"
	"gopkg.in/yaml.v3"
)

func executeRequest(
	ctx context.Context,
	ctr *app.Container,
	internalID int,
	req *PrefetchRequest,
	mapStore *sync.Map,
	hasDep bool,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctr.Logger.Debug(ctx, "executing request",
		logger.Value("id", req.ID), logger.Value("on", "executeRequest"))

	var targetReq PrefetchRequest
	if hasDep {
		ctr.Logger.Debug(ctx, "replace placeholders for dependent request",
			logger.Value("id", req.ID), logger.Value("on", "executeRequest"))
		output, err := yaml.Marshal(req)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %v", err)
		}
		placeholderRegex := regexp.MustCompile(`\<\.\.\<\s*(\w+)\s*\>\.\.\>`)

		replaced := placeholderRegex.ReplaceAllStringFunc(string(output), func(match string) string {
			key := placeholderRegex.FindStringSubmatch(match)[1]
			if v, exists := mapStore.Load(key); exists {
				return v.(string)
			}
			return match
		})

		if err := yaml.Unmarshal([]byte(replaced), &targetReq); err != nil {
			return fmt.Errorf("failed to unmarshal replaced request: %v", err)
		}

		ctr.Logger.Debug(ctx, "replaced request",
			logger.Value("id", req.ID), logger.Value("replaced", targetReq), logger.Value("on", "executeRequest"))
	}

	endType := req.EndpointType
	// INFO: close on factor.Factory(response handler), because only it will write to this channel
	queryTermChanWithBreak := make(chan queryreqbatch.TerminateType)
	queryType := queryreqbatch.NewQueryTypeFromString(endType)
	factor := queryreqbatch.TypeFactoryMap[queryType]

	writeFunc := func(
		ctx context.Context,
		ctr *app.Container,
		id int,
		data queryreqbatch.WriteData,
	) error {
		for _, replaceVar := range req.Vars {
			jmesPathQuery := replaceVar.JMESPath
			result, err := jmespath.Search(jmesPathQuery, data.RawData)
			if err != nil {
				ctr.Logger.Error(ctx, "failed to search jmespath",
					logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				return fmt.Errorf("failed to search jmespath: %v", err)
			}
			var ok bool
			if result != nil {
				if result, ok = result.(string); !ok {
					ctr.Logger.Warn(ctx, "result is not string",
						logger.Value("on", "runResponseHandler"), logger.Value("result", result))
				}
			}

			if ok {
				mapStore.Store(replaceVar.ID, result)
				ctr.Logger.Debug(ctx, "store value",
					logger.Value("id", replaceVar.ID), logger.Value("value", result), logger.Value("on", "runResponseHandler"))
			} else {
				ctr.Logger.Warn(ctx, "result is invalid",
					logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				switch replaceVar.OnError {
				case "ignore":
					ctr.Logger.Warn(ctx, "ignore error",
						logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				case "random":
					randomValue := utils.GenerateUniqueID()
					mapStore.Store(replaceVar.ID, randomValue)
					ctr.Logger.Warn(ctx, "random value",
						logger.Value("error", err), logger.Value("on", "runResponseHandler"), logger.Value("value", randomValue))
				case "empty":
					mapStore.Store(replaceVar.ID, "")
					ctr.Logger.Warn(ctx, "empty value",
						logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				case "cancel":
					ctr.Logger.Warn(ctx, "cancel value",
						logger.Value("error", err), logger.Value("on", "runResponseHandler"))
					return fmt.Errorf("cancel value")
				}
			}
		}

		return nil
	}
	exec, _, err := factor.Factory(
		ctx,
		ctr,
		internalID,
		&queryreqbatch.ValidatedQueryRequest{
			Endpoint:      req.EndpointType,
			Interval:      time.Microsecond * 1,
			AwaitPrevResp: false,
			QueryParam:    req.QueryParam,
			PathVariables: req.PathVariables,
			Break: queryreqbatch.NewSimpleValidatedQueryRequestBreak(
				time.Duration(10)*time.Minute,
				1,
				[]int{},
				struct {
					HasValue bool
					JMESPath string
				}{},
			),
			ExcludeStatusFilter: func(statusCode int) bool {
				return false
			},
		},
		queryTermChanWithBreak,
		ctr.AuthToken,
		ctr.Config.Web.QueryAPI.Url,
		writeFunc,
	)
	if err != nil {
		return fmt.Errorf("failed to create executor: %v", err)
	}

	err = exec.QueryExecute(ctx, ctr)
	if err != nil {
		ctr.Logger.Error(ctx, "failed to execute query",
			logger.Value("error", err), logger.Value("id", internalID))
		return fmt.Errorf("failed to execute query: %v", err)
	}

	termType := <-queryTermChanWithBreak
	ctr.Logger.Info(ctx, "Query End For Term For Break",
		logger.Value("QueryID", internalID), logger.Value("on", "executeRequest"))
	switch termType {
	case queryreqbatch.ByCount:
		ctr.Logger.Info(ctx, "Query End For Term By Count",
			logger.Value("QueryID", internalID), logger.Value("on", "executeRequest"))
	default:
		ctr.Logger.Error(ctx, "Query End For Error",
			logger.Value("QueryID", internalID), logger.Value("on", "executeRequest"))
		return fmt.Errorf("error because of termType: %v", queryreqbatch.TerminateType(termType).String())
	}
	return nil
}
