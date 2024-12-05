package batchtest

import (
	"fmt"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
)

func BatchTest(ctr *app.Container) error {

	filename, err := testprompt.PromptInput("File name: ")
	if err != nil {
		return fmt.Errorf("failed to get file name: %v", err)
	}

	globalStore := sync.Map{}
	ctx := ctr.Ctx

	if err := baseExecute(ctx, ctr, filename, &globalStore); err != nil {
		return fmt.Errorf("failed to execute batch test: %v", err)
	}

	return nil
}
