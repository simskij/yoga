package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DotsEncryptionConfig struct {
	Identity string `yaml:"identity"`
}

type DotsConfig struct {
	Encryption DotsEncryptionConfig `yaml:"encryption"`
}

type YoConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	Yo   YoConfig   `yaml:"yo"`
	Dots DotsConfig `yaml:"dots"`
}

func (c *Config) DotsPath() (string, error) {
	expanded, err := ExpandPath(c.Yo.Path)
	if err != nil {
		return "", err
	}
	return filepath.Join(expanded, "dots"), nil
}

func DefaultConfig() *Config {
	return &Config{
		Yo: YoConfig{
			Path: "~/.yofiles",
		},
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "yo", "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func Exists() (bool, error) {
	path, err := configPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func ExpandPath(path string) (string, error) {
	if len(path) > 1 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
