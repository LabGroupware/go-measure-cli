package massquery

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/executor"
	"github.com/LabGroupware/go-measure-tui/internal/api/response"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
)

type FindJobFactory struct{}

func (f FindJobFactory) factory(
	ctx context.Context,
	id int,
	interval time.Duration,
	responseWait bool,
	termChan chan<- struct{},
	authToken *auth.AuthToken,
	apiEndpoint string,
	outputFile *os.File,
) executor.RequestExecutor {
	var ok bool
	var err error
	var jobId string

	req := executor.GetJobReq{
		AuthToken:    authToken,
		BaseEndpoint: apiEndpoint,
	}

	for !ok {
		jobId, err = testprompt.PromptInput("ジョブIDを入力してください(random値を使用したい場合は`*`を入力)")
		if err != nil || jobId == "" {
			fmt.Println("ジョブIDに誤りがあります\n再度入力してください")
			continue
		}
		ok = true
		if jobId == "*" {
			jobId = testprompt.GenerateRandomString(10)
		}
		fmt.Printf("ジョブID: %s\n", jobId)
		req.Path.JobID = jobId
	}

	resChan := make(chan executor.ResponseContent[response.JobResponseDto])

	go func() {
		defer close(resChan)
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-resChan:
				_, _ = outputFile.WriteString(fmt.Sprintf("Response: %v\n", v))
				if v.HasSystemErr {
					fmt.Printf("リクエスト: %dの処理をシステムエラーにより中断します\n", id)
					termChan <- struct{}{}
				}
				// fmt.Printf("Response: %v\n", v)
			}
		}
	}()

	return executor.RequestContent[executor.GetJobReq, response.JobResponseDto]{
		Req:          req,
		Interval:     interval,
		ResponseWait: responseWait,
		ResChan:      resChan,
	}
}
