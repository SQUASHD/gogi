package command

import (
	"errors"
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
	"github.com/SQUASHD/gogi/internal/structs"
)

const (
	missingTemplate = "template name is required"
)

// commandCreate is the callback for the "add" command
// It adds a new template to the configuration
func (ctx *Context) commandCreate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(missingTemplate)
	}

	err := checkIfReservedWord(args[0])
	if err != nil {
		return err
	}

	err = ctx.handleCreate(args)
	if err != nil {
		return err
	}

	for _, arg := range args[1:] {
		switch arg {
		case "-e", "-edit":
			err = ctx.commandEdit(args)
			if err != nil {
				return err
			}
		case "-b", "-base":
			ctx.cfg.Base = args[0]
			if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
				return fmt.Errorf("could not save updated configuration: %w", err)
			}
		default:
			return fmt.Errorf("invalid flag '%s', expected -e, -edit, -b, or -base", arg)
		}
	}

	return nil
}

func (ctx *Context) handleCreate(args []string) error {
	if args[0] == "" {
		return fmt.Errorf(missingTemplate)
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
