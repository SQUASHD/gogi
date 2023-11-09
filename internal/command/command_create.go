package command

import (
	"errors"
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
	"github.com/SQUASHD/gogi/internal/structs"
)

// commandCreate is the callback for the "add" command
// It adds a new template to the configuration
func (ctx *Context) commandCreate(args []string) error {
	if len(args) == 0 || args[0] == "" {
		return fmt.Errorf("no template name provided to create")
	}
	name := args[0]
	_, err := config.FindTemplateByName(ctx.cfg, name)
	if err == nil {
		return errors.New("template already exists")
	}

	path := ctx.projectDir + "/" + name + ".gitignore"
	templ := &structs.Template{
		Name: name,
		Path: path,
	}
	ctx.cfg.Templates = append(ctx.cfg.Templates, *templ)

	if err = generator.CreateEmptyTemplateFile(ctx.projectDir, templ.Name); err != nil {
		return fmt.Errorf("could not create template file: %w", err)
	}

	if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}

	fmt.Printf("template '%s' created\n", name)
	return nil
}
