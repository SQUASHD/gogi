package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
)

// commandDelete is the callback for the "delete" command
func (ctx *Context) commandDelete(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided to delete")
	}

	name := args[0]
	deleted := false

	for i, templ := range ctx.cfg.Templates {
		if templ.Name == name {
			ctx.cfg.Templates = append(ctx.cfg.Templates[:i], ctx.cfg.Templates[i+1:]...)
			deleted = true
			fmt.Printf("template '%s' deleted\n", name)
			break
		}
	}
	if ctx.cfg.Base == name {
		ctx.cfg.Base = ""
		fmt.Println("base template deleted")
	}
	err := generator.DeleteTemplateFile(ctx.projectDir, name)
	if err != nil {
		return fmt.Errorf("could not delete template file: %w", err)
	}

	if !deleted {
		return fmt.Errorf("template '%s' not found", name)
	}

	if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}

	fmt.Printf("template '%s' deleted\n", name)
	return nil
}
