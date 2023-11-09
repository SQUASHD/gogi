package command

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"github.com/SQUASHD/gogi/internal/generator"
)

// HandleQuickGogi tries to create a .gitignore file based on the template
// designated as the base template
func (ctx *Context) HandleQuickGogi() error {
	baseTempl := ctx.cfg.Base
	if baseTempl == "" {
		return fmt.Errorf("no base template is set. try gogi base or gogi help")
	}
	templ, err := config.FindTemplateByName(ctx.cfg, baseTempl)
	if err != nil {
		return err
	}

	exists, err := generator.DoesGitignoreExist(ctx.cwd)
	if err != nil {
		return fmt.Errorf("error determining whether .gitignore exists")
	}
	if exists {
		fmt.Println("gogi with no argument is intended to run with no .gitignore file present")
		return fmt.Errorf("there's already a .gitignore file")
	}
	err = generator.GenerateGitignore(templ.Path, ctx.cwd)
	if err != nil {
		return err
	}
	fmt.Println("successfully added base .gitignore template")
	return nil
}
