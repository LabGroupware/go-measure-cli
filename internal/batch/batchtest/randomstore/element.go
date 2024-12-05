package randomstore

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"gopkg.in/yaml.v3"
)

type RandomElementValueGenerator struct {
	Key    string
	Values []interface{}
}

func (p *RandomElementValueGenerator) Generate(ctx context.Context, ctr *app.Container, store *sync.Map) error {
	if len(p.Values) == 0 {
		return fmt.Errorf("values is empty")
	}

	value := p.Values[rand.N(len(p.Values))]
	store.Store(p.Key, fmt.Sprintf("%v", value))
	return nil
}

type RandomStoreValueElementDataConfig struct {
	Key   string        `yaml:"key"`
	Type  string        `yaml:"type"`
	Value []interface{} `yaml:"value"`
}

func (p *RandomStoreValueElementDataConfig) Init(conf []byte) error {
	err := yaml.Unmarshal(conf, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return nil
}

func (p *RandomStoreValueElementDataConfig) GeneratorFactory(ctx context.Context, ctr *app.Container) (RadomGenerator, error) {
	return &RandomElementValueGenerator{
		Key:    p.Key,
		Values: p.Value,
	}, nil
}
