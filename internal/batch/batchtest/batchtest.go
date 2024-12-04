package batchtest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/massquerybatch"
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/prefetchbatch"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
	"gopkg.in/yaml.v3"
)

func BatchTest(ctr *app.Container) error {

	filename, err := testprompt.PromptInput("File name: ")
	if err != nil {
		return fmt.Errorf("failed to get file name: %v", err)
	}

	file, err := os.Open(filepath.Join(ctr.Config.Batch.Test.Path, filename))
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

	var reader io.Reader
	reader = file

	if conf.Prefetch.Enabled {
		var replacements map[string]string
		if replacements, err = prefetchbatch.PrefetchBatch(ctr, conf.Prefetch); err != nil {
			return fmt.Errorf("failed to execute prefetch: %v", err)
		}

		ctr.Logger.Debug(ctr.Ctx, "replacements set",
			logger.Value("replacements", replacements))

		var buffer bytes.Buffer
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			buffer.WriteString(scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		content := buffer.String()
		placeholderRegex := regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

		result := placeholderRegex.ReplaceAllStringFunc(content, func(match string) string {
			// プレースホルダー名を抽出
			key := placeholderRegex.FindStringSubmatch(match)[1]
			// マップにキーが存在する場合は置換、存在しない場合はそのまま
			if value, exists := replacements[key]; exists {
				return value
			}
			return match // キーが存在しない場合は元の文字列を返す
		})

		var yamlData map[string]interface{}

		if err := yaml.Unmarshal([]byte(result), &yamlData); err != nil {
			return fmt.Errorf("failed to parse as YAML: %w", err)
		}

		ctr.Logger.Debug(ctr.Ctx, "replaced content",
			logger.Value("content", yamlData))

		reader = bytes.NewReader([]byte(result))
	}

	switch conf.Type {
	case "MassQuery":
		var massQuery massquerybatch.MassQuery
		decoder := yaml.NewDecoder(reader)
		if err := decoder.Decode(&massQuery); err != nil {
			return fmt.Errorf("failed to decode yaml: %v", err)
		}
		if err := massquerybatch.MassQueryBatch(ctr, massQuery); err != nil {
			return fmt.Errorf("failed to execute mass query: %v", err)
		}
	case "WaitSaga":
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("unknown type")
	}

	return nil
}
