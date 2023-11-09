package command

import "fmt"

// commandHelp is the callback for the "help" command
// It displays the help message
func (ctx *Context) commandHelp(args []string) error {
	if len(args) > 0 {
		cmdName := resolveCommand(args[0])
		if cmd, ok := ctx.commands[cmdName]; ok {
			fmt.Printf("%s: %s\n", cmd.name, cmd.helpExample)
			return nil
		}
		return fmt.Errorf("unknown command: %s", cmdName)
	}
	var width int
	for _, cmd := range ctx.commands {
		if len(cmd.name) > width {
			width = len(cmd.name)
		}
	}
	for _, cmd := range ctx.commands {
		fmt.Printf("%*s: %s\n", width, cmd.name, cmd.description)
	}
	return nil
}
