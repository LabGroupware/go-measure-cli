package randomstore

import (
	"context"

	"github.com/LabGroupware/go-measure-tui/internal/app"
)

type RandomGenerator interface {
	Init(conf []byte) error
	GeneratorFactory(ctx context.Context, ctr *app.Container) (RadomGenerator, error)
}

var randomGeneratorFactoryMap = map[string]RandomGenerator{
	"element": &RandomStoreValueElementDataConfig{},
	"int":     &RandomStoreValueIntDataConfig{},
	"float":   &RandomStoreValueFloatDataConfig{},
	"string":  &RandomStoreValueStringDataConfig{},
	"bool":    &RandomStoreValueBoolDataConfig{},
	"uuid":    &RandomStoreValueUUIDDataConfig{},
}
