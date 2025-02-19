package app

import (
	"os"
	"time"

	"github.com/Tairascii/google-docs-documents/internal/app/usecase"
	"gopkg.in/yaml.v3"
)

const (
	configFilePath = "CONFIG_FILE_PATH"
)

type UseCase struct {
	Documents usecase.DocumentsUseCase
}

type DI struct {
	UseCase UseCase
}

type Config struct {
	Repo struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		DocumentsTable string `yaml:"documents_table"`
	} `yaml:"repo"`
	Server struct {
		Port    string `yaml:"port"`
		Timeout struct {
			Read  time.Duration `yaml:"read"`
			Write time.Duration `yaml:"write"`
			Idle  time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

func LoadConfigs() (*Config, error) {
	f, err := os.Open(os.Getenv(configFilePath))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
