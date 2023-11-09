package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
)

// commandRename handles renaming a template
func (ctx *Context) commandRename(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}

	err := checkIfReservedWord(args[1])
	if err != nil {
		return err
	}

	oldName := args[0]
	newName := args[1]
	templIdx, err := config.GetTemplateIndexByName(ctx.cfg, oldName)
	if err != nil {
		return fmt.Errorf("could not find template '%s'", oldName)
	}

	_, err = config.FindTemplateByName(ctx.cfg, newName)
	if err == nil {
		return fmt.Errorf("template '%s' already exists", newName)
	}

	if ctx.cfg.Base == oldName {
		ctx.cfg.Base = newName
		fmt.Printf("base template set to '%s'\n", newName)
	}

	newTemplatePath := generator.GenerateTemplatePath(ctx.projectDir, newName)
	ctx.cfg.Templates[templIdx].Name = newName
	if err := generator.RenameTemplateFile(ctx.projectDir, oldName, newName); err != nil {
		return fmt.Errorf("could not rename template file: %w", err)
	}
	ctx.cfg.Templates[templIdx].Path = newTemplatePath
	if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}
	fmt.Printf("template '%s' renamed to '%s'\n", oldName, newName)

	return nil
}
