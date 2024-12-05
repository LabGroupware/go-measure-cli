package executor

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type FindFileObjectReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		With []string
	}
	Path struct {
		FileObjectID string
	}
}

func (r FindFileObjectReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/file-objects"

	fileObjectID := r.Path.FileObjectID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, fileObjectID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to file object endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindFileObjectReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}

type GetFileObjectsReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		Limit           int
		Offset          int
		Cursor          string
		Pagination      string
		SortField       string
		SortOrder       string
		WithCount       bool
		HasBucketFilter bool
		FilterBucketIDs []string
		With            []string
	}
}

func (r GetFileObjectsReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/file-objects"

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
	queryParams.Set("has_bucket_filter", fmt.Sprint(r.Param.HasBucketFilter))
	for _, bucketID := range r.Param.FilterBucketIDs {
		queryParams.Add("filter_bucket_ids", bucketID)
	}
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to file objects endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "GetFileObjectsReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
