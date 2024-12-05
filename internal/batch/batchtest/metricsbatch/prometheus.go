package metricsbatch

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type PrometheusMetricsFetcher struct {
	req *http.Request
}

func (p *PrometheusMetricsFetcher) Fetch(ctx context.Context, ctr *app.Container) (any, error) {
	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// 	Transport: &utils.DelayedTransport{
	// 		Transport: http.DefaultTransport,
	// 		// Delay:     2 * time.Second,
	// 	},
	// }

	resp, err := http.Get(p.req.URL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Prometheus query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

type PrometheusMetricsBatchRequestConfig struct {
	ID       string                          `yaml:"id"`
	Type     string                          `yaml:"type"`
	URL      string                          `yaml:"url"`
	Query    string                          `yaml:"query"`
	Interval string                          `yaml:"interval"`
	Data     []MetricsBatchRequestDataConfig `yaml:"data"`
}

func (p *PrometheusMetricsBatchRequestConfig) FetcherFactory(ctx context.Context, ctr *app.Container) (MetricsFetcher, error) {

	baseURL := fmt.Sprintf("%s/api/v1/query", p.URL)
	fullURL, err := url.Parse(baseURL)

	queryParams := fullURL.Query()
	queryParams.Add("query", p.Query)
	fullURL.RawQuery = queryParams.Encode()

	ctr.Logger.Debug(ctx, "GET request to Prometheus query URL created",
		logger.Value("url", fullURL.String()), logger.Value("on", "PrometheusMetricsBatchRequestConfig.FetcherFactory"))

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	return &PrometheusMetricsFetcher{req: req}, nil
}

var _ MetricsFetcher = &PrometheusMetricsFetcher{}

var _ MetricsFetcherFactory = &PrometheusMetricsBatchRequestConfig{}
