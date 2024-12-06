package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type FindTaskReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		With []string
	}
	Path struct {
		TaskID string
	}
}

func (r FindTaskReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/tasks"

	taskID := r.Path.TaskID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, taskID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	ctr.Logger.Debug(ctx, "GET request to task endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindTaskReq.CreateRequest"))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type GetTasksReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		Limit                          int
		Offset                         int
		Cursor                         string
		Pagination                     string
		SortField                      string
		SortOrder                      string
		WithCount                      bool
		HasTeamFilter                  bool
		FilterTeamIDs                  []string
		HasStatusFilter                bool
		FilterStatuses                 []string
		HasChargeUserFilter            bool
		FilterChargeUserIDs            []string
		FilterStartDatetimeEarlierThan string
		FilterStartDatetimeLaterThan   string
		FilterDueDatetimeEarlierThan   string
		FilterDueDatetimeLaterThan     string
		HasFileObjectFilter            bool
		FilterFileObjectIDs            []string
		FileObjectFilterType           string
		With                           []string
	}
}

func (r GetTasksReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/teams"

	fullURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	queryParams.Set("limit", fmt.Sprint(r.Param.Limit))
	queryParams.Set("offset", fmt.Sprint(r.Param.Offset))
	queryParams.Set("cursor", r.Param.Cursor)
	queryParams.Set("pagination", r.Param.Pagination)
	queryParams.Set("sort_field", r.Param.SortField)
	queryParams.Set("sort_order", r.Param.SortOrder)
	queryParams.Set("with_count", fmt.Sprint(r.Param.WithCount))
	queryParams.Set("has_team_filter", fmt.Sprint(r.Param.HasTeamFilter))
	for _, teamID := range r.Param.FilterTeamIDs {
		queryParams.Add("filter_team_ids", teamID)
	}
	queryParams.Set("has_status_filter", fmt.Sprint(r.Param.HasStatusFilter))
	for _, status := range r.Param.FilterStatuses {
		queryParams.Add("filter_statuses", status)
	}
	queryParams.Set("has_charge_user_filter", fmt.Sprint(r.Param.HasChargeUserFilter))
	for _, chargeUserID := range r.Param.FilterChargeUserIDs {
		queryParams.Add("filter_charge_user_ids", chargeUserID)
	}
	queryParams.Set("filter_start_datetime_earlier_than", r.Param.FilterStartDatetimeEarlierThan)
	queryParams.Set("filter_start_datetime_later_than", r.Param.FilterStartDatetimeLaterThan)
	queryParams.Set("filter_due_datetime_earlier_than", r.Param.FilterDueDatetimeEarlierThan)
	queryParams.Set("filter_due_datetime_later_than", r.Param.FilterDueDatetimeLaterThan)
	queryParams.Set("has_file_object_filter", fmt.Sprint(r.Param.HasFileObjectFilter))
	for _, taskID := range r.Param.FilterFileObjectIDs {
		queryParams.Add("filter_file_object_ids", taskID)
	}
	queryParams.Set("file_object_filter_type", r.Param.FileObjectFilterType)
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to tasks endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "GetTasksReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type CreateTaskReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Body         any
}

func (r CreateTaskReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/tasks"

	fullURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctr.Logger.Debug(ctx, "POST request to task endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "CreateTaskReq.CreateRequest"))

	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fullURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type UpdateStatusTaskReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Path         struct {
		TaskID string
	}
	Body any
}

func (r UpdateStatusTaskReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/tasks"

	taskID := r.Path.TaskID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s/status", baseURL, taskID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctr.Logger.Debug(ctx, "PUT request to task status endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "UpdateStatusTaskReq.CreateRequest"))

	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fullURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
