package command

import (
	"fmt"
	"os"
)

// HandleCommand handles the incoming CLI arguments
func (ctx *Context) HandleCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("no command provided. try gogi help")
		return
	}

	cmdName := resolveCommand(args[0])
	if cmd, ok := ctx.commands[cmdName]; ok {
		if err := cmd.callback(ctx, args[1:]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: %s\n", cmdName)
	}
}

// Helper function to resolve command aliases
func resolveCommand(name string) string {
	if primaryName, exists := aliasMap[name]; exists {
		return primaryName
	}
	return name
}
