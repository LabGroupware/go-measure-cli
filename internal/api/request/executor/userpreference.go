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

type FindUserPreferenceReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Param        struct {
		With []string
	}
	Path struct {
		UserPreferenceID string
	}
}

func (r FindUserPreferenceReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/user-preferences"

	userPreferenceID := r.Path.UserPreferenceID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, userPreferenceID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := fullURL.Query()
	for _, with := range r.Param.With {
		queryParams.Add("with", with)
	}
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to user preference endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindUserPreferenceReq.CreateRequest"))

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

type UpdateUserPreferenceReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Path         struct {
		UserPreferenceID string
	}
	Body any
}

func (r UpdateUserPreferenceReq) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/user-preferences"

	userPreferenceID := r.Path.UserPreferenceID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, userPreferenceID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctr.Logger.Debug(ctx, "PUT request to userPreference endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "UpdateUserPreferenceReq.CreateRequest"))

	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
