package executor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/utils"
)

type ResponseContent[Res any] struct {
	Success         bool
	StartTime       time.Time
	EndTime         time.Time
	Res             Res
	ByteResponse    []byte
	ResponseTime    int64
	StatusCode      int
	ReqCreateHasErr bool
	ParseResHasErr  bool
	HasSystemErr    bool
	WithCountLimit  bool
}

type ExecReq interface {
	// CreateRequest creates the http.Request object for the query
	CreateRequest(ctx context.Context, ctr *app.Container) (*http.Request, error)
}

type RequestExecutor interface {
	RequestExecute(ctx context.Context, ctr *app.Container) error
}

type RequestCountLimit struct {
	Enabled bool
	Count   int
}

type RequestContent[Req ExecReq, Res any] struct {
	Req          Req
	Interval     time.Duration
	ResponseWait bool
	ResChan      chan<- ResponseContent[Res]
	CountLimit   RequestCountLimit
}

func (q RequestContent[Req, Res]) RequestExecute(
	ctx context.Context,
	ctr *app.Container,
) error {

	req, err := q.Req.CreateRequest(ctx, ctr)
	if err != nil {
		ctr.Logger.Error(ctx, "failed to create request",
			logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"))
		return err
	}

	go func() {
		// defer close(q.ResChan) // TODO: close channel

		ticker := time.NewTicker(q.Interval)
		defer ticker.Stop()
		var waitForResponse = q.ResponseWait
		var count int
		var countLimitOver bool
		chanForWait := make(chan struct{})
		defer close(chanForWait)

		for {
			select {
			case <-ctx.Done():
				ctr.Logger.Info(ctx, "request processing is interrupted due to context termination",
					logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
				return
			case <-ticker.C:
				if count > 0 && waitForResponse {
					select {
					case <-ctx.Done():
						ctr.Logger.Info(ctx, "request processing is interrupted due to context termination",
							logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))

						return
					case <-chanForWait:
					}
				}

				count++
				if q.CountLimit.Enabled && count >= q.CountLimit.Count {
					ctr.Logger.Info(ctx, "request processing is interrupted due to count limit",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					countLimitOver = true
				}

				reqClone := cloneRequest(req)

				go func(asyncReq *http.Request, countOver bool) {
					defer func() {
						if waitForResponse {
							chanForWait <- struct{}{}
						}
					}()

					client := &http.Client{
						Timeout: 10 * time.Minute,
						Transport: &utils.DelayedTransport{
							Transport: http.DefaultTransport,
							// Delay:     2 * time.Second,
						},
					}

					ctr.Logger.Debug(ctx, "sending request",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL), logger.Value("count", count))
					startTime := time.Now()
					resp, err := client.Do(asyncReq)
					endTime := time.Now()
					ctr.Logger.Debug(ctx, "received response",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL), logger.Value("count", count))
					if err != nil {
						ctr.Logger.Error(ctx, "response error",
							logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
						select {
						case <-ctx.Done():
							ctr.Logger.Info(ctx, "request processing is interrupted due to context termination",
								logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
							return
						case q.ResChan <- ResponseContent[Res]{
							Success:        false,
							StartTime:      startTime,
							EndTime:        endTime,
							ResponseTime:   endTime.Sub(startTime).Milliseconds(),
							HasSystemErr:   true,
							WithCountLimit: countOver,
						}: // do nothing
						}

						return
					}
					defer resp.Body.Close()

					statusCode := resp.StatusCode
					var response Res
					responseByte, err := io.ReadAll(resp.Body)
					if err != nil {
						ctr.Logger.Error(ctx, "failed to read response",
							logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
						select {
						case <-ctx.Done():
							ctr.Logger.Info(ctx, "request processing is interrupted due to context termination",
								logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))

							return
						case q.ResChan <- ResponseContent[Res]{
							Success:        false,
							Res:            response,
							StartTime:      startTime,
							EndTime:        endTime,
							ResponseTime:   endTime.Sub(startTime).Milliseconds(),
							StatusCode:     statusCode,
							ParseResHasErr: true,
							WithCountLimit: countOver,
						}: // do nothing
						}
						return
					}
					err = json.Unmarshal(responseByte, &response)
					if err != nil {
						ctr.Logger.Error(ctx, "failed to parse response",
							logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
						select {
						case <-ctx.Done():
							ctr.Logger.Info(ctx, "request processing is interrupted due to context termination",
								logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
							return
						case q.ResChan <- ResponseContent[Res]{
							Success:        false,
							Res:            response,
							ByteResponse:   responseByte,
							StartTime:      startTime,
							EndTime:        endTime,
							ResponseTime:   endTime.Sub(startTime).Milliseconds(),
							StatusCode:     statusCode,
							ParseResHasErr: true,
							WithCountLimit: countOver,
						}: // do nothing
						}
						return
					}

					ctr.Logger.Debug(ctx, "response OK",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					responseContent := ResponseContent[Res]{
						Success:        true,
						ByteResponse:   responseByte,
						Res:            response,
						StartTime:      startTime,
						EndTime:        endTime,
						ResponseTime:   endTime.Sub(startTime).Milliseconds(),
						StatusCode:     statusCode,
						WithCountLimit: countOver,
					}
					select {
					case q.ResChan <- responseContent:
					case <-ctx.Done():
						ctr.Logger.Info(ctx, "request processing is interrupted due to context termination",
							logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
						return
					}
				}(reqClone, countLimitOver)

				if countLimitOver {
					<-ctx.Done()

					ctr.Logger.Info(ctx, "request processing is interrupted due to count limit",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					return
				}
			}
		}
	}()

	return nil
}

func cloneRequest(req *http.Request) *http.Request {
	clone := req.Clone(req.Context())
	if req.Body != nil {
		body, _ := req.GetBody()
		clone.Body = body
	}
	return clone
}
