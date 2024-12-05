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

type FindTeamReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		With []string
	}
	Path struct {
		TeamID string
	}
}

func (r FindTeamReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/teams"

	teamID := r.Path.TeamID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, teamID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to team endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindTeamReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type GetTeamsReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		Limit                 int
		Offset                int
		Cursor                string
		Pagination            string
		SortField             string
		SortOrder             string
		WithCount             bool
		HasIsDefaultFilter    bool
		FilterIsDefault       bool
		HasOrganizationFilter bool
		FilterOrganizationIDs []string
		HasUserFilter         bool
		FilterUserIDs         []string
		UserFilterType        string
		With                  []string
	}
}

func (r GetTeamsReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
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
	queryParams.Set("has_is_default_filter", fmt.Sprint(r.Param.HasIsDefaultFilter))
	queryParams.Set("filter_is_default", fmt.Sprint(r.Param.FilterIsDefault))
	queryParams.Set("has_organization_filter", fmt.Sprint(r.Param.HasOrganizationFilter))
	for _, organizationID := range r.Param.FilterOrganizationIDs {
		queryParams.Add("filter_organization_ids", organizationID)
	}
	queryParams.Set("has_user_filter", fmt.Sprint(r.Param.HasUserFilter))
	for _, userID := range r.Param.FilterUserIDs {
		queryParams.Add("filter_user_ids", userID)
	}
	queryParams.Set("user_filter_type", r.Param.UserFilterType)
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to teams endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "GetTeamsReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type CreateTeamReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Body         any
}

func (r CreateTeamReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/teams"

	fullURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctr.Logger.Debug(ctx, "POST request to team endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "CreateTeamReq.CreateRequest"))

	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fullURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type AddUsersTeamReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Path         struct {
		TeamID string
	}
	Body any
}

func (r AddUsersTeamReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/teams"

	teamID := r.Path.TeamID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s/users", baseURL, teamID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctr.Logger.Debug(ctx, "POST request to team users endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "AddUsersTeamReq.CreateRequest"))

	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fullURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
