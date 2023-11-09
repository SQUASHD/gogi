package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
	"os"
)

// commandGenerate is the callback for the "generate" command
// It generates a .gitignore file from the given template
func (ctx *Context) commandGenerate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided")
	}

	name := args[0]
	force := false
	if len(args) > 1 {
		force = args[1] == "--force" || args[1] == "-f"
	}

	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return err
	}

	if err := generator.CheckWhetherTemplateExists(templ.Path); err != nil {
		return err
	}

	exists, err := generator.DoesGitignoreExist(ctx.cwd)
	if err != nil {
		return err
	}

	if exists && !force {
		confirmMsg := "A .gitignore file already exists. Do you want to overwrite it?"
		confirmed, err := ConfirmAction(confirmMsg, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println(OperationCancelledString)
			return nil
		}
	}

	if err := generator.GenerateGitignore(templ.Path, ctx.cwd); err != nil {
		return err
	}

	fmt.Printf("Generated .gitignore file from template '%s'\n", name)
	return nil
}
