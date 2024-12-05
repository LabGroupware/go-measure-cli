package metricsbatch

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

type PrometheusMetricsFetcher struct {
	req      *http.Request
	interval time.Duration
}

func (r PrometheusMetricsFetcher) CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error) {
	return r.req, nil
}

func (p *PrometheusMetricsFetcher) Fetch(
	ctx context.Context,
	ctr *app.Container,
) (any, chan<- struct{}, error) {

	// INFO: close on executor, because only it will write to this channel
	resChan := make(chan queryreq.ResponseContent[any])
	// INFO: close on factor.Factory(response handler), because only it will write to this channel
	termChan := make(chan TermType)

	go func() {
		for {
			select {
			case <-ctx.Done():
				termChan <- TermTypeContext
				return
			}
		}
	}()

	req := queryreq.RequestContent[PrometheusMetricsFetcher, any]{
		Req:          *p,
		Interval:     p.interval,
		ResponseWait: false,
		ResChan:      resChan,
		CountLimit:   queryreq.RequestCountLimit{},
	}

	execTerm, err := req.QueryExecute(ctx, ctr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	go func() {
		defer close(termChan)

		<-termChan
		execTerm <- struct{}{}
		ctr.Logger.Info(ctx, "Prometheus Query End For Term",
			logger.Value("on", "PrometheusMetricsFetcher.Fetch"))
	}()

	return nil, execTerm, nil

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
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

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

	var interval time.Duration

	if p.Interval != "" {
		interval, err = time.ParseDuration(p.Interval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse interval: %w", err)
		}
	} else {
		return nil, fmt.Errorf("interval is required")
	}

	return &PrometheusMetricsFetcher{
		req:      req,
		interval: interval,
	}, nil
}

var _ MetricsFetcher = &PrometheusMetricsFetcher{}

var _ MetricsFetcherFactory = &PrometheusMetricsBatchRequestConfig{}
