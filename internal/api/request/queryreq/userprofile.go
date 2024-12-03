package queryreq

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type FindUserReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		With []string
	}
	Path struct {
		UserID string
	}
}

func (r FindUserReq) CreateRequest(ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/user-profiles"

	userID := r.Path.UserID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, userID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctr.Ctx, "GET request to user profile endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindUserReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type GetUsersReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		Limit      int
		Offset     int
		Cursor     string
		Pagination string
		SortField  string
		SortOrder  string
		WithCount  bool
		With       []string
	}
}

func (r GetUsersReq) CreateRequest(ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/user-profiles"

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
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctr.Ctx, "GET request to user profiles endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "GetUsersReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}