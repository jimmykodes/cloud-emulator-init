package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jimmykodes/cloud-emulator-init/internal/aws"
	"github.com/jimmykodes/cloud-emulator-init/internal/gcp"
)

type Config struct {
	AWS *aws.Config `yaml:"aws,omitempty"`
	GCP *gcp.Config `yaml:"gcp,omitempty"`
}

func New(file string) (*Config, error) {
	var conf Config
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &conf, yaml.NewDecoder(f).Decode(&conf)
}
