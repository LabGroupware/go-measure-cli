package queryreqbatch

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/google/uuid"
	"github.com/jmespath/go-jmespath"
)

type writeData struct {
	Success          bool
	SendDatetime     string
	ReceivedDatetime string
	Count            int
	ResponseTime     int
	StatusCode       string
	Data             any
}

func (d writeData) ToSlice() []string {
	return []string{
		strconv.FormatBool(d.Success),
		d.SendDatetime,
		d.ReceivedDatetime,
		strconv.Itoa(d.Count),
		strconv.Itoa(d.ResponseTime),
		d.StatusCode,
		fmt.Sprintf("%v", d.Data),
	}
}

type writeSendData struct {
	uid       uuid.UUID
	writeData []string
}

func runResponseHandler[Res any](
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	uidChan <-chan uuid.UUID,
	resChan <-chan queryreq.ResponseContent[Res],
	writeChan chan<- writeSendData,
) {
	var count int
	var timeout <-chan time.Time
	if request.Break.Time > 0 {
		timeout = time.After(request.Break.Time)
	} else {
		timeout = make(chan time.Time)
	}
	sentUid := make(map[uuid.UUID]struct{})
	for {
		if request.Break.Count.HasValue && count == request.Break.Count.Value {
			ctr.Logger.Info(ctr.Ctx, "Term Condition: Count",
				logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
			for len(sentUid) > 0 {
				uid := <-uidChan
				delete(sentUid, uid)
			}
			termChan <- struct{}{}
			return
		}
		count++
		select {
		case uid := <-uidChan:
			delete(sentUid, uid)
			count--
		case <-timeout:
			ctr.Logger.Info(ctr.Ctx, "Term Condition: Time",
				logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
			for len(sentUid) > 0 {
				uid := <-uidChan
				delete(sentUid, uid)
			}
			termChan <- struct{}{}
		case <-ctr.Ctx.Done():
			ctr.Logger.Info(ctr.Ctx, "Term Condition: Context Done",
				logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
			return
		case v := <-resChan:
			var mustWrite bool
			var response interface{}
			err := json.Unmarshal(v.ByteResponse, &response)
			if err != nil {
				ctr.Logger.Error(ctr.Ctx, "The response is not a valid JSON",
					logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				for len(sentUid) > 0 {
					uid := <-uidChan
					delete(sentUid, uid)
				}
				termChan <- struct{}{}
			}

			if request.DataOutputFilter.HasValue {
				jmesPathQuery := request.DataOutputFilter.JMESPath
				fmt.Println(jmesPathQuery)
				result, err := jmespath.Search(jmesPathQuery, response)
				if err != nil {
					ctr.Logger.Error(ctr.Ctx, "failed to search jmespath",
						logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				}
				if result != nil {
					if v, ok := result.(bool); ok {
						if v {

							mustWrite = true
						}
					} else {
						ctr.Logger.Warn(ctr.Ctx, "The result of the jmespath query is not a boolean",
							logger.Value("on", "runResponseHandler"))
					}
				}
			} else {
				mustWrite = true
			}

			if request.ExcludeStatusFilter(v.StatusCode) {
				ctr.Logger.Info(ctr.Ctx, "Status output filter found",
					logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
				mustWrite = false
			}

			if mustWrite {
				var data any = v.Res
				if request.DataOutput.HasValue {
					jmesPathQuery := request.DataOutput.JMESPath
					result, err := jmespath.Search(jmesPathQuery, response)
					if err != nil {
						ctr.Logger.Error(ctr.Ctx, "failed to search jmespath",
							logger.Value("error", err), logger.Value("on", "runResponseHandler"))
					}
					data = result
				}
				uid := uuid.New()
				writeData := writeData{
					Success:          v.Success,
					SendDatetime:     v.StartTime.Format(time.RFC3339),
					ReceivedDatetime: v.EndTime.Format(time.RFC3339),
					Count:            count,
					ResponseTime:     int(v.ResponseTime),
					StatusCode:       strconv.Itoa(v.StatusCode),
					Data:             data,
				}
				sentUid[uid] = struct{}{}
				writeChan <- writeSendData{
					uid:       uid,
					writeData: writeData.ToSlice(),
				}
			}
			if v.ReqCreateHasErr {
				ctr.Logger.Warn(ctr.Ctx, "Term Condition: Request Creation Error",
					logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
				for len(sentUid) > 0 {
					uid := <-uidChan
					delete(sentUid, uid)
				}
				termChan <- struct{}{}
				return
			}
			if v.HasSystemErr {
				if request.Break.SysError {
					ctr.Logger.Warn(ctr.Ctx, "Term Condition: System Error",
						logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
					for len(sentUid) > 0 {
						uid := <-uidChan
						delete(sentUid, uid)
					}
					termChan <- struct{}{}
					return
				} else {
					ctr.Logger.Warn(ctr.Ctx, "System error occurred",
						logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
				}
			}
			if v.ParseResHasErr {
				if request.Break.ParseError {
					ctr.Logger.Warn(ctr.Ctx, "Term Condition: Response Parse Error",
						logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
					for len(sentUid) > 0 {
						uid := <-uidChan
						delete(sentUid, uid)
					}
					termChan <- struct{}{}
					return
				} else {
					ctr.Logger.Warn(ctr.Ctx, "Parse error occurred",
						logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
				}
			}
			if request.Break.ResponseBody.HasValue {
				jmesPathQuery := request.Break.ResponseBody.JMESPath
				result, err := jmespath.Search(jmesPathQuery, response)
				if err != nil {
					ctr.Logger.Error(ctr.Ctx, "failed to search jmespath",
						logger.Value("error", err), logger.Value("on", "runResponseHandler"))
				}
				if result != nil {
					ctr.Logger.Info(ctr.Ctx, "Term Condition: Response Body",
						logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
					for len(sentUid) > 0 {
						uid := <-uidChan
						delete(sentUid, uid)
					}
					termChan <- struct{}{}
					return
				}
			}
			if request.Break.StatusCodeMatcher(v.StatusCode) {
				ctr.Logger.Info(ctr.Ctx, "Term Condition: Status Code",
					logger.Value("id", id), logger.Value("count", count), logger.Value("on", "runResponseHandler"))
				for len(sentUid) > 0 {
					uid := <-uidChan
					delete(sentUid, uid)
				}
				termChan <- struct{}{}
				return
			}
		}
	}
}

func runAsyncProcessing[Res any](
	ctr *app.Container,
	id int,
	request *ValidatedQueryRequest,
	termChan chan<- struct{},
	resChan chan queryreq.ResponseContent[Res],
	outputFile *os.File,
) {
	writeChan := make(chan writeSendData)
	wroteUidChan := make(chan uuid.UUID)
	go func() {
		defer close(resChan)
		runResponseHandler(ctr, id, request, termChan, wroteUidChan, resChan, writeChan)
	}()

	go func() {
		defer close(writeChan)
		for {
			select {
			case d := <-writeChan:
				writer := csv.NewWriter(outputFile)
				ctr.Logger.Debug(ctr.Ctx, "Writing data to csv",
					logger.Value("id", id), logger.Value("data", d), logger.Value("on", "runAsyncProcessing"))
				if err := writer.Write(d.writeData); err != nil {
					ctr.Logger.Error(ctr.Ctx, "failed to write data to csv",
						logger.Value("error", err), logger.Value("on", "runAsyncProcessing"))
				}
				writer.Flush()
				wroteUidChan <- d.uid
			case <-ctr.Ctx.Done():
				return
			}
		}
	}()
}
