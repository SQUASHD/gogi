package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
)

// commandBase handles setting the base template or callsback the edit
// command to edit the base template based on the user flags
func (ctx *Context) commandBase(args []string) error {
	baseName := ctx.cfg.Base
	if len(args) == 0 && baseName == "" {
		return fmt.Errorf("no base template set")
	} else if len(args) == 0 {
		fmt.Printf("Your current base file is template: %v\n", baseName)
		return nil
	}

	name := args[0]
	if name == "" {
		return fmt.Errorf("no template name provided")
	}
	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return fmt.Errorf("no template with that name exists")
	}
	ctx.cfg.Base = templ.Name
	if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
		return err
	}
	fmt.Printf("base template set to '%s'\n", name)
	return nil
}
