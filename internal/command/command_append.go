package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
)

// commandAppend is the callback for the "append" command
// It appends a template to an existing gitignore file
func (ctx *Context) commandAppend(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided to append")
	}

	name := args[0]
	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return fmt.Errorf("could not find template '%s'", name)
	}

	exists, err := generator.DoesGitignoreExist(ctx.cwd)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("couldn't find a gitignore file to append to.")
	}

	if err := generator.AppendTemplate(ctx.cwd, templ.Path); err != nil {
		return err
	}

	fmt.Printf("appended template '%s' to gitignore file\n", name)
	return nil
}
