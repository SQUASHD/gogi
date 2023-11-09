package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
)

// commandGenerate is the callback for the "generate" command
// It generates a gitignore file from the given template
func (ctx *Context) commandGenerate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided")
	}
	name := args[0]
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
	if exists {
		return fmt.Errorf("a gitignore file already exists")
	}
	if err := generator.GenerateGitignore(templ.Path, ctx.cwd); err != nil {
		return err
	}
	fmt.Printf("generated gitignore file from template '%s'\n", name)
	return nil
}
