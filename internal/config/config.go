package config

import (
	"errors"
	"fmt"
	goconfig "github.com/SQUASHD/go-config/config"
	"github.com/SQUASHD/gogi/internal/structs"
	"os"
)

var ErrTemplateNotFound = errors.New("template not found")

var Dir = os.Getenv("HOME") + "/.config/gogi"
var configPath = Dir + "/gogi.json"

func InitConfig() error {
	var cfg structs.TemplateConfig
	if err := goconfig.InitConfig(configPath, cfg); err != nil {
		return fmt.Errorf("could not initialize configuration at %s: %w", configPath, err)
	}
	return nil
}

func LoadConfig() (*structs.TemplateConfig, error) {
	var cfg structs.TemplateConfig
	if err := goconfig.LoadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("could not load config file. Try gogi init")
	}
	return &cfg, nil
}

func SaveConfig(cfg *structs.TemplateConfig) error {
	if err := goconfig.SaveConfig(configPath, cfg); err != nil {
		return fmt.Errorf("could not save configuration to %s: %w", configPath, err)
	}
	return nil
}

func FindTemplateByName(cfg *structs.TemplateConfig, name string) (*structs.Template, error) {
	for _, tmpl := range cfg.Templates {
		if tmpl.Name == name {
			return &tmpl, nil
		}
	}
	return nil, ErrTemplateNotFound
}
