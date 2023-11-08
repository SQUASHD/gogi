package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/SQUASHD/gogi/internal/commands"
	"github.com/SQUASHD/gogi/internal/generator"
	"github.com/SQUASHD/gogi/pkg/config"
)

func RunCli(args []string) {
	cwd := os.Getenv("PWD")
	if len(args) < 2 {
		if err := handleQuickGogi(cwd); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	switch args[1] {
	case "init":
		if err := config.InitConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Configuration initialized successfully.")
	default:
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ctx, err := commands.NewCommandContext(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		args := sanitizeArgs(args[1:])
		ctx.HandleCommand(args)
	}
}

// handleQuickGogi tries to create a .gitignore file based on the template
// designated as the base template
func handleQuickGogi(cwd string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	baseTempl := cfg.Base
	if baseTempl == "" {
		return fmt.Errorf("no base template is set. try gogi base or gogi help")
	}
	templ, err := config.FindTemplateByName(cfg, baseTempl)
	if err != nil {
		return err
	}

	exists, err := generator.DoesGitignoreExist(cwd)
	if err != nil {
		return fmt.Errorf("error determining whether .gitignore exists")
	}
	if exists {
		fmt.Println("gogi with no argument is intended to run with no .gitignore file present")
		return fmt.Errorf("there's already a .gitignore file")
	}
	err = generator.GenerateGitignore(templ.Path, cwd)
	if err != nil {
		return err
	}
	fmt.Println("successfully added base .gitignore template")
	return nil
}

func sanitizeArgs(args []string) []string {
	sanitizedArgs := []string{}
	for _, arg := range args {
		sanitizedArgs = append(sanitizedArgs, strings.ToLower(arg))
	}
	return sanitizedArgs
}
