package queryreq

import (
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

func (r FindUserPreferenceReq) CreateRequest(ctr *app.Container) (*http.Request, error) {
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

	ctr.Logger.Debug(ctr.Ctx, "GET request to user preference endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "FindUserPreferenceReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
