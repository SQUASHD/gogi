package structs

import "github.com/SQUASHD/go-config/config"

type TemplateConfig struct {
	Editor    string     `json:"editor"`
	Base      string     `json:"base"`
	Templates []Template `json:"templates"`
}

type Template struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (c TemplateConfig) Default() config.Config {
	return TemplateConfig{
		Editor:    "nano",
		Templates: []Template{},
	}
}
