package queryreq

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type FindOrganizationReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		With []string
	}
	Path struct {
		OrganizationID string
	}
}

func (r FindOrganizationReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/organizations"

	organizationID := r.Path.OrganizationID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, organizationID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to organization endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindOrganizationReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type GetOrganizationsReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		Limit          int
		Offset         int
		Cursor         string
		Pagination     string
		SortField      string
		SortOrder      string
		WithCount      bool
		HasOwnerFilter bool
		FilterOwnerIDs []string
		HasPlanFilter  bool
		FilterPlans    []string
		HasUserFilter  bool
		FilterUserIDs  []string
		UserFilterType string
		With           []string
	}
}

func (r GetOrganizationsReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/organizations"

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
	queryParams.Set("has_owner_filter", fmt.Sprint(r.Param.HasOwnerFilter))
	for _, ownerID := range r.Param.FilterOwnerIDs {
		queryParams.Add("filter_owner_ids", ownerID)
	}
	queryParams.Set("has_plan_filter", fmt.Sprint(r.Param.HasPlanFilter))
	for _, plan := range r.Param.FilterPlans {
		queryParams.Add("filter_plans", plan)
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

	ctr.Logger.Debug(ctx, "GET request to organizations endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "GetOrganizationsReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
