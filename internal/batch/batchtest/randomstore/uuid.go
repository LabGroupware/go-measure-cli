package randomstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type RandomUUIDValueGenerator struct {
	Key string
}

func (p *RandomUUIDValueGenerator) Generate(ctx context.Context, ctr *app.Container, store *sync.Map) error {
	values := uuid.New()
	store.Store(p.Key, fmt.Sprintf("%v", values.String()))
	return nil
}

type RandomStoreValueUUIDDataConfig struct {
	Key  string `yaml:"key"`
	Type string `yaml:"type"`
}

func (p *RandomStoreValueUUIDDataConfig) Init(conf []byte) error {
	err := yaml.Unmarshal(conf, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return nil
}

func (p *RandomStoreValueUUIDDataConfig) GeneratorFactory(ctx context.Context, ctr *app.Container) (RadomGenerator, error) {
	return &RandomUUIDValueGenerator{
		Key: p.Key,
	}, nil
}
