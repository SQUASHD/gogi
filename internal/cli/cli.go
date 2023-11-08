package cli

import (
	"fmt"
	"github.com/SQUASHD/gogi/internal/command"
	"github.com/SQUASHD/gogi/internal/config"
	"os"
	"strings"
)

func RunCli(args []string) {
	if len(args) > 1 && args[1] == "init" {
		if err := config.InitConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Configuration initialized successfully.")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := command.NewCommandContext(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) == 1 {
		if err := ctx.HandleQuickGogi(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	args = sanitizeArgs(args[1:])
	ctx.HandleCommand(args)
}

func sanitizeArgs(args []string) []string {
	sanitizedArgs := []string{}
	for _, arg := range args {
		sanitizedArgs = append(sanitizedArgs, strings.ToLower(arg))
	}
	return sanitizedArgs
}
