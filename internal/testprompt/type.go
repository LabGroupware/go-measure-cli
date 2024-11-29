package testprompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type TestPromptID int

const (
	_ TestPromptID = iota
	TestPromptMassiveQuery
	TestPromptWaitSaga
	TestPromptUrgeOnConsistentAfterStartSaga
	TestPromptUrgeOnConsistentAfterEndSaga
)

type TestPromptData struct {
	ID          TestPromptID
	Name        string
	Description string
}

func TestTypeSelection() (TestPromptID, error) {
	types := []TestPromptData{
		{ID: TestPromptMassiveQuery, Name: "Massive Query", Description: "大量のデータをクエリを送信する"},
		{ID: TestPromptWaitSaga, Name: "Wait Saga", Description: "サーガの処理をWebSocketsで待機して, 通知を受け取る"},
		{ID: TestPromptUrgeOnConsistentAfterStartSaga, Name: "Urge on Consistent After Start Saga", Description: "サーガの処理が開始してから, 整合性を確認するまでクエリを送信し続ける"},
		{ID: TestPromptUrgeOnConsistentAfterEndSaga, Name: "Urge on Consistent After End Saga", Description: "サーガの処理が完了してから, 整合性を確認するまでクエリを送信し続ける"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U00002714 {{ .ID | cyan }} {{ .Name | green }}",
		Inactive: "  {{ .ID | cyan }} ({{ .Name | white }})",
		Selected: "\U00002714 {{ .ID | red | cyan }}",
		Details: `
--------- タスク詳細 ----------
{{ "ID:" | faint }}	{{ .ID }}
{{ "名前:" | faint }} {{ .Name }}
{{ "説明:" | faint }} {{ .Description }}`,
	}

	prompt := promptui.Select{
		Label:     "テストタイプを選択してください",
		Items:     types,
		Templates: templates,
		Size:      4,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("選択に失敗しました: %v\n", err)
		return 0, err
	}

	selectedTask := types[index]
	fmt.Printf("選択したタスク: %s\n", selectedTask.Name)

	return selectedTask.ID, nil
}
