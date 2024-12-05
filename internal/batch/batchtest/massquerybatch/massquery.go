package massquerybatch

import "github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/queryreqbatch"

type MassQueryData struct {
	Requests []MassQueryOneData `yaml:"requests"`
}

type MassQueryOneData struct {
	queryreqbatch.QueryRequest `yaml:",inline"`
	SuccessBreak               []string `yaml:"successBreak"`
}

type MassQuery struct {
	Type   string          `yaml:"type"`
	Data   MassQueryData   `yaml:"data"`
	Output BatchTestOutput `yaml:"output"`
}

type BatchTestOutput struct {
	Enabled bool `yaml:"enabled"`
}
