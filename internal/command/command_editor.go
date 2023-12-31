package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
)

// commandEditor is the callback for the "editor" command
// It sets the editor to use for editing templates
func (ctx *Context) commandEditor(args []string) error {
	if len(args) == 0 && ctx.cfg.Editor == "" {
		fmt.Printf("you have not set an editor")
		return nil
	}
	if len(args) == 0 {
		fmt.Printf("editor set to '%s'\n", ctx.cfg.Editor)
		return nil
	}
	name := args[0]

	if name == "" {
		return fmt.Errorf("no editor name provided")
	}
	ctx.cfg.Editor = name
	if err := config.SaveConfig(ctx.cfg, ctx.configPath); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}
	fmt.Printf("editor set to '%s'\n", name)
	return nil
}
