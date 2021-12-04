package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/miracl/conflate"
)

type AppConfig struct {
	data map[string]json.RawMessage
}

func NewAppConfig() *AppConfig {
	merge := conflate.New()
	merge.AddFiles("config.json")

	for _, env := range os.Environ() {
		merge.AddGo(env) //todo formatting
	}
	//Load other sources
	//e.g. os.Args, vault, secrets manager etc.

	merged, _ := merge.MarshalJSON()
	var data map[string]json.RawMessage
	_ = json.Unmarshal(merged, &data)

	return &AppConfig{data}
}

func (cfg *AppConfig) Value(conf interface{}) error {
	if reflect.TypeOf(conf).Kind() != reflect.Ptr {
		return errors.New("config is not a pointer type")
	}
	configName := reflect.TypeOf(conf).Elem().Name()

	raw, ok := cfg.data[configName]
	if !ok {
		return nil
	}

	return json.Unmarshal(raw, conf)
}

// FooServiceConfiguration ...
type FooServiceConfiguration struct {
	Host      string
	APIKey    string
	TimeoutMs int
	UseBeta   bool
}

// NewFooServiceConfiguration ...
func NewFooServiceConfiguration(ac *AppConfig) (FooServiceConfiguration, error) {
	cfg := &FooServiceConfiguration{}
	err := ac.Value(cfg)
	return *cfg, err
}
