package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
	"os"
)

// commandDelete is the callback for the "delete" command
func (ctx *Context) commandDelete(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided to delete")
	}

	name := args[0]
	forced := len(args) > 1 && (args[1] == "--f" || args[1] == "--force")

	if len(args) > 2 || (len(args) == 2 && !forced) {
		return fmt.Errorf("invalid arguments provided")
	}

	var confirmationPrompt string
	if !forced {
		if ctx.cfg.Base == name {
			confirmationPrompt = fmt.Sprintf("Are you sure?\n\nTemplate '%s' is currently the base template.", name)
		} else {
			confirmationPrompt = fmt.Sprintf("Are you sure you want to delete template '%s'?", name)
		}
		confirmed, err := ConfirmAction(confirmationPrompt, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}
	}

	if err := ctx.deleteTemplate(name); err != nil {
		return err
	}

	fmt.Printf("template '%s' deleted\n", name)
	return nil
}

// deleteTemplate handles the deletion of a template from the configuration and the file system.
func (ctx *Context) deleteTemplate(name string) error {
	deleted := false
	for i, templ := range ctx.cfg.Templates {
		if templ.Name == name {
			ctx.cfg.Templates = append(ctx.cfg.Templates[:i], ctx.cfg.Templates[i+1:]...)
			deleted = true
			break
		}
	}

	if !deleted {
		return fmt.Errorf("template '%s' not found", name)
	}

	if ctx.cfg.Base == name {
		ctx.cfg.Base = ""
		fmt.Println("base template deleted")
	}

	if err := generator.DeleteTemplateFile(ctx.projectDir, name); err != nil {
		return fmt.Errorf("could not delete template file: %w", err)
	}

	if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}

	return nil
}
