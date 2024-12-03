package massquerybatch

import "github.com/LabGroupware/go-measure-tui/internal/batch/batchtest/queryreqbatch"

type MassQueryData struct {
	Requests []queryreqbatch.QueryRequest `yaml:"requests"`
}

type MassQuery struct {
	Type string        `yaml:"type"`
	Data MassQueryData `yaml:"data"`
}
