package massquery

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
	"github.com/manifoldco/promptui"
)

func MassiveQueryPrompt(apiEndpoint string, authToken *auth.AuthToken) {
	var ok bool
	var concurrentCount int
	var err error
	ctx := context.Background()

	for !ok {
		concurrentCount, err = testprompt.PromptNumber("いくつ同時に並行してクエリを叩くか")
		if err != nil {
			fmt.Printf("入力に誤りがあります\n 再度入力してください: %v\n", err)
			continue
		}
		ok = true
	}

	if concurrentCount <= 0 {
		return
	}

	queryExecutors := make([]*MassiveQueryThreadExecutor, concurrentCount)

	timestamp := time.Now().Format("20060102_150405")
	dirPath := fmt.Sprintf("./bench/%s", timestamp)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Printf("ディレクトリ作成に失敗しました: %v\n", err)
		return
	}

	for i := 0; i < concurrentCount; i++ {
		logFilePath := fmt.Sprintf("%s/logfile_%d.txt", dirPath, i+1)
		file, err := os.Create(logFilePath)
		if err != nil {
			fmt.Printf("ログファイル作成に失敗しました: %v\n", err)
			return
		}
		queryExecutors[i] = &MassiveQueryThreadExecutor{
			ID:         i + 1,
			outputFile: file,
		}
	}

	for i := 0; i < concurrentCount; i++ {
		fmt.Printf("Query %d の設定を入力してください\n", i+1)

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U00002714 {{ .ID | cyan }} {{ .Name | green }}",
			Inactive: "  {{ .ID | cyan }} ({{ .Name | white }})",
			Selected: "\U00002714 {{ .ID | red | cyan }}",
			Details: `
	--------- クエリタイプ詳細 ----------
	{{ "ID:" | faint }}	{{ .ID }}
	{{ "名前:" | faint }} {{ .Name }}
	{{ "説明:" | faint }} {{ .Description }}`,
		}

		prompt := promptui.Select{
			Label:     "テストタイプを選択してください",
			Items:     QueryTypes,
			Templates: templates,
			Size:      8,
		}

		var index int
		var queryType QueryType
		var err error

		ok = false

		for !ok {
			index, _, err = prompt.Run()
			if err != nil {
				fmt.Printf("入力に誤りがあります\n 再度入力してください: %v\n", err)
				continue
			} else {
				ok = true
				queryType = QueryTypes[index].ID
			}
		}

		ok = false
		var interval time.Duration

		for !ok {
			interval, err = testprompt.PromptMillisecond("クエリの時間間隔（ms）を入力してください")
			if err != nil {
				fmt.Printf("入力に誤りがあります\n 再度入力してください: %v\n", err)
				continue
			}
			ok = true
		}

		ok = false
		var responseWait bool

		for !ok {
			responseWait, err = testprompt.PromptBool("ResponseWait（レスポンスを待つか？）")
			if err != nil {
				fmt.Printf("入力に誤りがあります\n 再度入力してください: %v\n", err)
				continue
			}
			ok = true
		}

		factory := typeFactoryMap[queryType]
		termChan := make(chan struct{})
		executor := factory.factory(ctx, i+1, interval, responseWait, termChan, authToken, apiEndpoint, queryExecutors[i].outputFile)

		fmt.Printf("Query %d の設定\n", i+1)
		fmt.Printf("[QueryType]: %v\n", QueryTypes[index].Name)
		fmt.Printf("[Interval]: %v\n", interval)
		fmt.Printf("[ResponseWait]: %v\n", responseWait)

		queryExecutors[i].RequestExecutor = executor
		queryExecutors[i].TermChan = termChan
	}

	massiveQueryExecutor := NewMassiveQueryExecutorWithThreads(queryExecutors)
	defer massiveQueryExecutor.Close(ctx)

	if err := massiveQueryExecutor.Execute(ctx); err != nil {
		fmt.Printf("クエリ実行中にエラーが発生しました: %v\n", err)
	}
}
