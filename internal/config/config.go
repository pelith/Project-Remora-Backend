package config

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config[T any] struct {
	Env   string
	Name  string `mapstructure:"name" structs:"name"`
	Pprof bool   `mapstructure:"pprof" structs:"pprof"`
	Log   struct {
		Level slog.Level `mapstructure:"level" structs:"level"`
	} `mapstructure:"log" structs:"log"`
	AppConfig T `mapstructure:"app_config" structs:"app_config"`
}

func Load[T any](env string) (*Config[T], error) {
	return LoadFromDir[T](env, "./config")
}

func LoadFromDir[T any](env string, configDir string) (*Config[T], error) {
	var cfg Config[T]

	cfgMap := structs.Map(cfg)
	flatCfgMap := make(map[string]interface{})
	flatten("", cfgMap, flatCfgMap)

	v := viper.New()

	for key := range flatCfgMap {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))

		err := v.BindEnv(key, envKey)
		if err != nil {
			return nil, fmt.Errorf("bind env: %w", err)
		}
	}

	v.SetConfigFile(fmt.Sprintf("%s/base.yaml", configDir))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	v.SetConfigFile(fmt.Sprintf("%s/%s.yaml", configDir, env))

	err := v.MergeInConfig()
	if err != nil {
		return nil, fmt.Errorf("merge config: %w", err)
	}

	hooks := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		logLevelHookFunction,
		uuidHookFunction,
	)

	err = v.Unmarshal(&cfg, viper.DecodeHook(hooks))
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.Env = env

	return &cfg, nil
}

func flatten(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}

	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flatten(prefix+k, child, dest)
		case []interface{}:
			for i := range child {
				dest[prefix+k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			dest[prefix+k] = v
		}
	}
}

func logLevelHookFunction(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(slog.Level(0)) {
		return data, nil
	}

	if f.Kind() == reflect.String {
		var level slog.Level

		levelStr, ok := data.(string)
		if !ok {
			return nil, ErrInvalidLogLevel
		}

		err := level.UnmarshalText([]byte(levelStr))
		if err != nil {
			return nil, fmt.Errorf("unmarshal log level: %w", err)
		}

		return level, nil
	}

	return nil, &UnsupportedTypeError{
		Type: f.Kind().String(),
	}
}

func uuidHookFunction(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(uuid.UUID{}) {
		return data, nil
	}

	if f.Kind() == reflect.String {
		uuidStr, ok := data.(string)
		if !ok {
			return nil, ErrInvalidUUID
		}

		id, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("parse UUID: %w", err)
		}

		return id, nil
	}

	return nil, &UnsupportedTypeError{
		Type: f.Kind().String(),
	}
}
