package command

import "fmt"

// commandList is the callback for the "list" command
// It lists all the templates in the configuration
func (ctx *Context) commandList(args []string) error {
	if len(ctx.cfg.Templates) == 0 {
		fmt.Println("you don't have any templates!")
		fmt.Println("try gogi create template-name to create a new one")
		return nil
	}

	fmt.Println("Available templates:")
	for _, templ := range ctx.cfg.Templates {
		fmt.Println("- " + templ.Name)
	}

	return nil
}
