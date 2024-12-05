package randomstore

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"gopkg.in/yaml.v3"
)

type RandomBoolValueGenerator struct {
	Key string
}

func (p *RandomBoolValueGenerator) Generate(ctx context.Context, ctr *app.Container, store *sync.Map) error {
	values := []bool{true, false}
	store.Store(p.Key, fmt.Sprintf("%v", values[rand.N(len(values))]))
	return nil
}

type RandomStoreValueBoolDataConfig struct {
	Key  string `yaml:"key"`
	Type string `yaml:"type"`
}

func (p *RandomStoreValueBoolDataConfig) Init(conf []byte) error {
	err := yaml.Unmarshal(conf, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return nil
}

func (p *RandomStoreValueBoolDataConfig) GeneratorFactory(ctx context.Context, ctr *app.Container) (RadomGenerator, error) {
	return &RandomBoolValueGenerator{
		Key: p.Key,
	}, nil
}
