package core

import (
	"os"

	"github.com/BurntSushi/toml"
)

type GopkgToml struct {
	Name         string            `toml:"name"`
	Dependencies map[string]string `toml:"dependencies"`
}

func LoadToml(path string) (*GopkgToml, error) {
	var cfg GopkgToml
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Dependencies == nil {
		cfg.Dependencies = make(map[string]string)
	}

	return &cfg, nil
}

func SaveToml(path string, cfg *GopkgToml) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(cfg)
}
