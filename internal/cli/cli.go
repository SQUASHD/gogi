package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/SQUASHD/gogi/internal/command"
	"github.com/SQUASHD/gogi/internal/config"
)

//var projectDir = os.Getenv("HOME") + "/.config/gogi"
//var configPath = projectDir + "/config.json"

var projectDir = "/Users/hjartland/repos/cli-utils/quick-gi/testing"
var configPath = projectDir + "/config.json"

func RunCli(args []string) {
	args = sanitizeArgs(args)

	if len(args) > 1 && args[1] == "init" {
		if err := config.InitConfig(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Configuration initialized successfully.")
		return
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cwd := os.Getenv("PWD")
	ctx, err := command.NewCommandContext(cfg, cwd, projectDir, configPath)
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

	args = args[1:]
	ctx.HandleCommand(args)
}

func sanitizeArgs(args []string) []string {
	sanitizedArgs := []string{}
	for _, arg := range args {
		sanitizedArgs = append(sanitizedArgs, strings.ToLower(arg))
	}
	return sanitizedArgs
}
