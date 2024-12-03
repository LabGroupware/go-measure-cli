package batchtest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/massquerybatch"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
	"gopkg.in/yaml.v3"
)

func BatchTest(container *app.Container) error {

	filename, err := testprompt.PromptInput("File name: ")
	if err != nil {
		return fmt.Errorf("failed to get file name: %v", err)
	}

	file, err := os.Open(filepath.Join(container.Config.Batch.Test.Path, filename))
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var conf BatchTestType
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		return fmt.Errorf("failed to decode yaml: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	switch conf.Type {
	case "MassQuery":
		var massQuery massquerybatch.MassQuery
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&massQuery); err != nil {
			return fmt.Errorf("failed to decode yaml: %v", err)
		}
		if err := massquerybatch.MassQueryBatch(container, massQuery); err != nil {
			return fmt.Errorf("failed to execute mass query: %v", err)
		}
	case "WaitSaga":
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("unknown type")
	}

	return nil
}
