package queryreq

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type GetJobReq struct {
	BaseEndpoint string
	AuthToken    *auth.AuthToken
	Path         struct {
		JobID string
	}
}

func (r GetJobReq) CreateRequest(ctr *app.Container) (*http.Request, error) {
	baseURL := r.BaseEndpoint + "/jobs"

	jobID := r.Path.JobID
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, jobID))
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctr.Logger.Debug(ctr.Ctx, "GET request to job endpoint URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "GetJobReq.CreateRequest"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	r.AuthToken.SetAuthHeader(req)

	return req, nil
}
