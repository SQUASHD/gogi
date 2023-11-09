package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
)

// commandEdit is the callback for the "edit" command
// It edits an existing template
func (ctx *Context) commandEdit(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided")
	}
	name := args[0]
	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return err
	}

	err = openTemplateInEditor(ctx.cfg.Editor, templ.Path)
	if err != nil {
		return err
	}

	return nil
}
