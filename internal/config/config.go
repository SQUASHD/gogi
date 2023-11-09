package config

import (
	"errors"
	"fmt"
	goconfig "github.com/SQUASHD/go-config/config"
	"github.com/SQUASHD/gogi/internal/structs"
)

var ErrTemplateNotFound = errors.New("template not found")

func InitConfig(configPath string) error {
	var cfg structs.TemplateConfig
	if err := goconfig.InitConfig(configPath, cfg); err != nil {
		return fmt.Errorf("could not initialize configuration at %s: %w", configPath, err)
	}
	return nil
}

func LoadConfig(configPath string) (*structs.TemplateConfig, error) {
	var cfg structs.TemplateConfig
	if err := goconfig.LoadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("could not load config file. Try gogi init")
	}
	return &cfg, nil
}

func SaveConfig(cfg *structs.TemplateConfig, configPath string) error {
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

func GetTemplateIndexByName(cfg *structs.TemplateConfig, name string) (int, error) {
	for i, tmpl := range cfg.Templates {
		if tmpl.Name == name {
			return i, nil
		}
	}
	return -1, ErrTemplateNotFound
}

func AddTemplate(cfg *structs.TemplateConfig, tmpl structs.Template) error {
	cfg.Templates = append(cfg.Templates, tmpl)
	return nil
}

func UpdateTemplate(cfg *structs.TemplateConfig, tmpl structs.Template, index int) error {
	cfg.Templates[index] = tmpl
	return nil
}
