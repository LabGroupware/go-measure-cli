package queryreqbatch

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type QueryRequest struct {
	EndpointType      string              `yaml:"endpointType"`
	Interval          string              `yaml:"interval"`
	AwaitPrevResponse bool                `yaml:"awaitPrevResponse"`
	QueryParam        map[string][]string `yaml:"queryParam"`
	PathVariables     map[string]string   `yaml:"pathVariables"`
	Break             struct {
		Time       *string `yaml:"time"`
		Count      *int    `yaml:"count"`
		SysError   *bool   `yaml:"sysError"`
		ParseError *bool   `yaml:"parseError"`
		WriteError *bool   `yaml:"writeError"`
		StatusCode *struct {
			Op    *string `yaml:"op"`
			Value *string `yaml:"value"`
		} `yaml:"statusCode"`
		ResponseBody *struct {
			JMESPath *string `yaml:"jmesPath"`
		} `yaml:"responseBody"`
	} `yaml:"break"`
	DataOutput *struct {
		JMESPath *string `yaml:"jmesPath"`
	} `yaml:"dataOutput"`
	ExcludeStatusFilter *struct {
		Op    *string `yaml:"op"`
		Value *string `yaml:"value"`
	} `yaml:"excludeStatusFilter"`
	DataOutputFilter *struct {
		JMESPath *string `yaml:"jmesPath"`
	} `yaml:"dataOutputFilter"`
}

func ValidateQueryReq(ctx context.Context, ctr *app.Container, req QueryRequest, validated *ValidatedQueryRequest) error {
	queryType := NewQueryTypeFromString(req.EndpointType)
	if queryType == 0 {
		return fmt.Errorf("invalid query type: %s", req.EndpointType)
	}
	validated.Endpoint = req.EndpointType
	interval, err := time.ParseDuration(req.Interval)
	if err != nil {
		return fmt.Errorf("failed to parse interval: %w", err)
	}
	validated.Interval = interval
	validated.AwaitPrevResp = req.AwaitPrevResponse
	validated.QueryParam = req.QueryParam
	validated.PathVariables = req.PathVariables
	if req.Break.Time != nil {
		breakTime, err := time.ParseDuration(*req.Break.Time)
		if err != nil {
			return fmt.Errorf("failed to parse break time: %w", err)
		}
		validated.Break.Time = breakTime
	}
	if req.Break.Count != nil {
		validated.Break.Count.Enabled = true
		validated.Break.Count.Count = *req.Break.Count
	}
	if req.Break.SysError != nil {
		validated.Break.SysError = *req.Break.SysError
	}
	if req.Break.ParseError != nil {
		validated.Break.ParseError = *req.Break.ParseError
	}
	if req.Break.WriteError != nil {
		validated.Break.WriteError = *req.Break.WriteError
	}
	op := "none"
	value := "0"
	if req.Break.StatusCode != nil {
		if req.Break.StatusCode.Op == nil || req.Break.StatusCode.Value == nil {
			return fmt.Errorf("status code operator and value must be set")
		}
		op = *req.Break.StatusCode.Op
		value = *req.Break.StatusCode.Value
	}
	statusCodeMatcher, err := statusCodeMatherFactory(ctx, ctr, op, value)
	if err != nil {
		return fmt.Errorf("failed to create status code matcher: %w", err)
	}
	validated.Break.StatusCodeMatcher = statusCodeMatcher
	if req.Break.ResponseBody != nil {
		if req.Break.ResponseBody.JMESPath != nil {
			validated.Break.ResponseBody.HasValue = true
			validated.Break.ResponseBody.JMESPath = *req.Break.ResponseBody.JMESPath
		}
	}
	if req.DataOutput != nil {
		if req.DataOutput.JMESPath != nil {
			validated.DataOutput.HasValue = true
			validated.DataOutput.JMESPath = *req.DataOutput.JMESPath
		}
	}
	op = "none"
	value = "0"
	if req.ExcludeStatusFilter != nil {
		if req.ExcludeStatusFilter.Op == nil || req.ExcludeStatusFilter.Value == nil {
			return fmt.Errorf("status code operator and value must be set")
		}
		op = *req.ExcludeStatusFilter.Op
		value = *req.ExcludeStatusFilter.Value
	}
	statusCodeMatcher, err = statusCodeMatherFactory(ctx, ctr, op, value)
	if err != nil {
		return fmt.Errorf("failed to create status code matcher: %w", err)
	}
	validated.ExcludeStatusFilter = statusCodeMatcher
	if req.DataOutputFilter != nil {
		if req.DataOutputFilter.JMESPath != nil {
			validated.DataOutputFilter.HasValue = true
			validated.DataOutputFilter.JMESPath = *req.DataOutputFilter.JMESPath
		}
	}

	return nil
}

type StatusCodeMatcher func(statusCode int) bool

func statusCodeMatherFactory(ctx context.Context, ctr *app.Container, op string, value string) (StatusCodeMatcher, error) {
	switch op {
	case "none":
		return func(statusCode int) bool {
			return false
		}, nil
	case "eq":
		statusCodeInt, err := strconv.Atoi(value)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode == statusCodeInt
		}, nil
	case "ne":
		statusCodeInt, err := strconv.Atoi(value)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode != statusCodeInt
		}, nil
	case "lt":
		statusCodeInt, err := strconv.Atoi(value)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode < statusCodeInt
		}, nil
	case "le":
		statusCodeInt, err := strconv.Atoi(value)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode <= statusCodeInt
		}, nil
	case "gt":
		statusCodeInt, err := strconv.Atoi(value)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode > statusCodeInt
		}, nil
	case "ge":
		statusCodeInt, err := strconv.Atoi(value)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode >= statusCodeInt
		}, nil
	case "in":
		statusCodeStrings := strings.Split(value, ",")
		codes := make([]int, len(statusCodeStrings))
		for _, v := range statusCodeStrings {
			statusCodeInt, err := strconv.Atoi(v)
			if err != nil {
				ctr.Logger.Error(ctx, "failed to convert status code to int",
					logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
				return nil, err
			}
			codes = append(codes, statusCodeInt)
		}
		return func(statusCode int) bool {
			for _, v := range codes {
				if statusCode == v {
					return true
				}
			}
			return false
		}, nil
	case "nin":
		statusCodeStrings := strings.Split(value, ",")
		codes := make([]int, len(statusCodeStrings))
		for _, v := range statusCodeStrings {
			statusCodeInt, err := strconv.Atoi(v)
			if err != nil {
				ctr.Logger.Error(ctx, "failed to convert status code to int",
					logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
				return nil, err
			}
			codes = append(codes, statusCodeInt)
		}
		return func(statusCode int) bool {
			for _, v := range codes {
				if statusCode == v {
					return false
				}
			}
			return true
		}, nil
	case "between":
		statusCodeStrings := strings.Split(value, ",")
		if len(statusCodeStrings) != 2 {
			return nil, fmt.Errorf("between operator must have 2 values")
		}
		min, err := strconv.Atoi(statusCodeStrings[0])
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		max, err := strconv.Atoi(statusCodeStrings[1])
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode >= min && statusCode <= max
		}, nil
	case "notBetween":
		statusCodeStrings := strings.Split(value, ",")
		if len(statusCodeStrings) != 2 {
			return nil, fmt.Errorf("between operator must have 2 values")
		}
		min, err := strconv.Atoi(statusCodeStrings[0])
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		max, err := strconv.Atoi(statusCodeStrings[1])
		if err != nil {
			ctr.Logger.Error(ctx, "failed to convert status code to int",
				logger.Value("error", err), logger.Value("on", "statusCodeMatherFactory"))
			return nil, err
		}
		return func(statusCode int) bool {
			return statusCode < min || statusCode > max
		}, nil
	case "regex":
		preRequireRegex := regexp.MustCompile(value)
		return func(statusCode int) bool {
			return preRequireRegex.MatchString(strconv.Itoa(statusCode))
		}, nil
	default:
		ctr.Logger.Error(ctx, "unknown operator",
			logger.Value("operator", op), logger.Value("on", "statusCodeMatherFactory"))
		return nil, fmt.Errorf("unknown operator: %s", op)
	}
}

type ValidatedQueryRequestBreak struct {
	Time              time.Duration
	Count             queryreq.RequestCountLimit
	SysError          bool
	ParseError        bool
	WriteError        bool
	StatusCodeMatcher StatusCodeMatcher
	ResponseBody      struct {
		HasValue bool
		JMESPath string
	}
}

func NewSimpleValidatedQueryRequestBreak(
	timeout time.Duration,
	count int,
	statusesForTerm []int,
	responseBody struct {
		HasValue bool
		JMESPath string
	},
) ValidatedQueryRequestBreak {
	return ValidatedQueryRequestBreak{
		Time: timeout,
		Count: queryreq.RequestCountLimit{
			Enabled: true,
			Count:   count,
		},
		SysError:   true,
		ParseError: true,
		WriteError: true,
		StatusCodeMatcher: func(statusCode int) bool {
			for _, v := range statusesForTerm {
				if statusCode == v {
					return true
				}
			}
			return false
		},
		ResponseBody: responseBody,
	}
}

type ValidatedQueryRequestDataOutput struct {
	HasValue bool
	JMESPath string
}

type ValidatedQueryRequestDataOutputFilter struct {
	HasValue bool
	JMESPath string
}

type ValidatedQueryRequest struct {
	Endpoint            string
	Interval            time.Duration
	AwaitPrevResp       bool
	QueryParam          map[string][]string
	PathVariables       map[string]string
	Break               ValidatedQueryRequestBreak
	DataOutput          ValidatedQueryRequestDataOutput
	ExcludeStatusFilter StatusCodeMatcher
	DataOutputFilter    ValidatedQueryRequestDataOutputFilter
}
