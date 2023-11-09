package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
	"os"
)

// HandleQuickGogi tries to create a .gitignore file based on the template
// designated as the base template
func (ctx *Context) HandleQuickGogi() error {
	baseTempl := ctx.cfg.Base
	if baseTempl == "" {
		return fmt.Errorf("no base template is set. try 'gogi base' or 'gogi help'")
	}
	templ, err := config.FindTemplateByName(ctx.cfg, baseTempl)
	if err != nil {
		return err
	}

	exists, err := generator.DoesGitignoreExist(ctx.cwd)
	if err != nil {
		return fmt.Errorf("error determining whether .gitignore exists: %v", err)
	}

	if exists {
		confirmationPrompt := "A .gitignore file already exists.\nOverwrite?"
		confirmed, err := ConfirmAction(confirmationPrompt, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println(OperationCancelledString)
			return nil
		}
	}

	err = generator.GenerateGitignore(templ.Path, ctx.cwd)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created .gitignore template from base.")
	return nil
}
