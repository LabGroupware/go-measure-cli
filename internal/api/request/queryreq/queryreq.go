package queryreq

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
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
}

type QueryReq interface {
	// CreateRequest creates the http.Request object for the query
	CreateRequest(ctr *app.Container) (*http.Request, error)
}

type QueryExecutor interface {
	QueryExecute(ctr *app.Container) (chan<- struct{}, error)
}

type RequestContent[Req QueryReq, Res any] struct {
	Req          Req
	Interval     time.Duration
	ResponseWait bool
	ResChan      chan<- ResponseContent[Res]
}

func (q RequestContent[Req, Res]) QueryExecute(
	ctr *app.Container,
) (chan<- struct{}, error) {

	terminateChan := make(chan struct{})

	req, err := q.Req.CreateRequest(ctr)
	if err != nil {
		ctr.Logger.Error(ctr.Ctx, "failed to create request",
			logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"))
		close(terminateChan)
		return nil, err
	}

	go func() {
		defer close(terminateChan)

		ticker := time.NewTicker(q.Interval)
		defer ticker.Stop()

		var previousResponse *ResponseContent[Res]
		var waitForResponse = q.ResponseWait

		for {
			select {
			case <-terminateChan:
				ctr.Logger.Info(ctr.Ctx, "request processing is interrupted due to termination",
					logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
				return
			case <-ctr.Ctx.Done():
				ctr.Logger.Info(ctr.Ctx, "request processing is interrupted due to context termination",
					logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
				return
			case <-ticker.C:
				if waitForResponse && previousResponse != nil {
					select {
					case <-terminateChan:
						ctr.Logger.Info(ctr.Ctx, "request processing is interrupted due to termination",
							logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
						return
					case <-ctr.Ctx.Done():
						ctr.Logger.Info(ctr.Ctx, "request processing is interrupted due to context termination",
							logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
						return
					case q.ResChan <- *previousResponse:
						ctr.Logger.Debug(ctr.Ctx, "previous response OK",
							logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					}
				}

				reqClone := cloneRequest(req)

				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				ctr.Logger.Debug(ctr.Ctx, "sending request",
					logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
				startTime := time.Now()
				resp, err := client.Do(reqClone)
				endTime := time.Now()
				ctr.Logger.Debug(ctr.Ctx, "received response",
					logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
				if err != nil {
					ctr.Logger.Error(ctr.Ctx, "failed to send request",
						logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					q.ResChan <- ResponseContent[Res]{
						Success:      false,
						StartTime:    startTime,
						EndTime:      endTime,
						ResponseTime: endTime.Sub(startTime).Milliseconds(),
						HasSystemErr: true,
					}
					continue
				}
				defer resp.Body.Close()

				statusCode := resp.StatusCode
				var response Res
				responseByte, err := io.ReadAll(resp.Body)
				if err != nil {
					ctr.Logger.Error(ctr.Ctx, "failed to read response",
						logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					q.ResChan <- ResponseContent[Res]{
						Success:        false,
						Res:            response,
						StartTime:      startTime,
						EndTime:        endTime,
						ResponseTime:   endTime.Sub(startTime).Milliseconds(),
						StatusCode:     statusCode,
						ParseResHasErr: true,
					}
					continue
				}
				err = json.Unmarshal(responseByte, &response)
				if err != nil {
					ctr.Logger.Error(ctr.Ctx, "failed to parse response",
						logger.Value("error", err), logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					q.ResChan <- ResponseContent[Res]{
						Success:        false,
						Res:            response,
						ByteResponse:   responseByte,
						StartTime:      startTime,
						EndTime:        endTime,
						ResponseTime:   endTime.Sub(startTime).Milliseconds(),
						StatusCode:     statusCode,
						ParseResHasErr: true,
					}
					continue
				}

				ctr.Logger.Debug(ctr.Ctx, "response OK",
					logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
				responseContent := ResponseContent[Res]{
					Success:      true,
					ByteResponse: responseByte,
					Res:          response,
					StartTime:    startTime,
					EndTime:      endTime,
					ResponseTime: endTime.Sub(startTime).Milliseconds(),
					StatusCode:   statusCode,
				}

				select {
				case q.ResChan <- responseContent:
					previousResponse = &responseContent
				case <-terminateChan:
					ctr.Logger.Info(ctr.Ctx, "request processing is interrupted due to termination",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					return
				case <-ctr.Ctx.Done():
					ctr.Logger.Info(ctr.Ctx, "request processing is interrupted due to context termination",
						logger.Value("on", "RequestContent.QueryExecute"), logger.Value("url", req.URL))
					return
				}
			}
		}
	}()

	return terminateChan, nil
}

func cloneRequest(req *http.Request) *http.Request {
	clone := req.Clone(req.Context())
	if req.Body != nil {
		body, _ := req.GetBody()
		clone.Body = body
	}
	return clone
}
