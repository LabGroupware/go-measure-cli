package randomstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"gopkg.in/yaml.v3"
)

type ConstantValueGenerator struct {
	Key   string
	Value string
}

func (p *ConstantValueGenerator) Generate(ctx context.Context, ctr *app.Container, store *sync.Map) error {
	store.Store(p.Key, p.Value)
	return nil
}

type RandomStoreValueConstantDataConfig struct {
	Key   string `yaml:"key"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

func (p *RandomStoreValueConstantDataConfig) Init(conf []byte) error {
	err := yaml.Unmarshal(conf, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return nil
}

func (p *RandomStoreValueConstantDataConfig) GeneratorFactory(ctx context.Context, ctr *app.Container) (RadomGenerator, error) {
	return &ConstantValueGenerator{
		Key:   p.Key,
		Value: p.Value,
	}, nil
}
