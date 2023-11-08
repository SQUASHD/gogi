package command

import (
	"errors"
	"fmt"
	"github.com/SQUASHD/gogi/internal/config"
	"os"
	"os/exec"

	"github.com/SQUASHD/gogi/internal/generator"
	"github.com/SQUASHD/gogi/internal/structs"
)

var aliasMap = map[string]string{
	"h": "help",
	"c": "create",
	"l": "list",
	"g": "generate",
	"e": "edit",
	"d": "delete",
	"a": "append",
	"b": "base",
}

// Context holds the state and provides methods to execute CLI commands
type Context struct {
	cfg        *structs.TemplateConfig
	cwd        string
	commands   map[string]cliCommand
	projectDir string
}

// cliCommand represents a command in the CLI
type cliCommand struct {
	name        string
	description string
	fullCommand string
	callback    func(*Context, []string) error
}

// NewCommandContext initializes a new command context with the given configuration
// and current working directory
func NewCommandContext(cfg *structs.TemplateConfig) (*Context, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get current working directory: %w", err)
	}
	ctx := &Context{
		cfg:        cfg,
		cwd:        cwd,
		projectDir: config.Dir,
	}
	ctx.commands = ctx.getCommands()
	return ctx, nil
}

// getCommands returns a map of commands with their corresponding callback functions
func (ctx *Context) getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"create": {
			name:        "create",
			description: "Create a new template",
			fullCommand: "gogi create template-name",
			callback:    (*Context).commandCreate,
		},
		"delete": {
			name:        "delete",
			description: "Delete an existing gitignore alias",
			fullCommand: "gogi delete template-name",
			callback:    (*Context).commandDelete,
		},
		"list": {
			name:        "list",
			description: "List all the templates",
			fullCommand: "gogi list",
			callback:    (*Context).commandList,
		},
		"generate": {
			name:        "generate",
			description: "Generate a gitignore file from the given template",
			fullCommand: "gogi generate template-name",
			callback:    (*Context).commandGenerate,
		},
		"edit": {
			name:        "edit",
			description: "Edit an existing template",
			fullCommand: "gogi edit template-name",
			callback:    (*Context).commandEdit,
		},
		"append": {
			name:        "append",
			description: "Append a template to an existing gitignore file",
			fullCommand: "gogi append template-name",
			callback:    (*Context).commandAppend,
		},
		"help": {
			name:        "help",
			description: "Display help message, or help for a specific command",
			fullCommand: "gogi help [command]",
			callback:    (*Context).commandHelp,
		},
		"editor": {
			name:        "editor",
			description: "Set the editor to use for editing templates",
			fullCommand: "gogi editor editor-name",
			callback:    (*Context).commandEditor,
		},
		"base": {
			name:        "base",
			description: "set the base template that you call with gogi with no args",
			fullCommand: "gogi base template-name",
			callback:    (*Context).commandBase,
		},
		"alias": {
			name:        "alias",
			description: "show the list of avaiable command aliases",
			fullCommand: "gogi alias",
			callback:    (*Context).commandAlias,
		},
	}
}

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

// commandDelete is the callback for the "delete" command
func (ctx *Context) commandDelete(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided to delete")
	}

	name := args[0]
	deleted := false

	for i, templ := range ctx.cfg.Templates {
		if templ.Name == name {
			ctx.cfg.Templates = append(ctx.cfg.Templates[:i], ctx.cfg.Templates[i+1:]...)
			deleted = true
			fmt.Printf("template '%s' deleted\n", name)
			break
		}
	}
	if ctx.cfg.Base == name {
		ctx.cfg.Base = ""
		fmt.Println("base template deleted")
	}
	err := generator.DeleteTemplateFile(ctx.projectDir, name)
	if err != nil {
		return fmt.Errorf("could not delete template file: %w", err)
	}

	if !deleted {
		return fmt.Errorf("template '%s' not found", name)
	}

	if err := config.SaveConfig(ctx.cfg); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}

	fmt.Printf("template '%s' deleted\n", name)
	return nil
}

// commandCreate is the callback for the "add" command
// It adds a new template to the configuration
func (ctx *Context) commandCreate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided to create")
	}

	name := args[0]
	_, err := config.FindTemplateByName(ctx.cfg, name)
	if err == nil {
		return errors.New("template already exists")
	}

	path := ctx.projectDir + "/" + name + ".gitignore"
	templ := &structs.Template{
		Name: name,
		Path: path,
	}
	ctx.cfg.Templates = append(ctx.cfg.Templates, *templ)

	if err = generator.CreateEmptyTemplateFile(ctx.projectDir, templ.Name); err != nil {
		return fmt.Errorf("could not create template file: %w", err)
	}

	if err := config.SaveConfig(ctx.cfg); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}

	fmt.Printf("template '%s' created\n", name)
	return nil
}

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

// commandEdit is the callback for the "edit" command
// It edits an existing template
func (ctx *Context) commandEdit(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided")
	}
	name := args[0]
	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return err
	}

	err = openTemplateInEditor(ctx.cfg.Editor, templ.Path)
	if err != nil {
		return err
	}

	return nil
}

// commandAppend is the callback for the "append" command
// It appends a template to an existing gitignore file
func (ctx *Context) commandAppend(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template name provided to append")
	}

	name := args[0]
	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return fmt.Errorf("could not find template '%s'", name)
	}

	// check if .gitignore exists
	exists, err := generator.DoesGitignoreExist(ctx.cwd)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("couldn't find a gitignore file to append to.")
	}

	if err := generator.AppendTemplate(ctx.cwd, templ.Path); err != nil {
		return err
	}

	fmt.Printf("appended template '%s' to gitignore file\n", name)
	return nil
}

// commandHelp is the callback for the "help" command
// It displays the help message
func (ctx *Context) commandHelp(args []string) error {
	if len(args) > 0 {
		cmdName := resolveCommand(args[0])
		if cmd, ok := ctx.commands[cmdName]; ok {
			fmt.Printf("%s: %s\n", cmd.name, cmd.fullCommand)
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

// commandEditor is the callback for the "editor" command
// It sets the editor to use for editing templates
func (ctx *Context) commandEditor(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no editor name provided")
	}
	name := args[0]

	if name == "" {
		return fmt.Errorf("no editor name provided")
	}
	ctx.cfg.Editor = name
	if err := config.SaveConfig(ctx.cfg); err != nil {
		return fmt.Errorf("could not save updated configuration: %w", err)
	}
	fmt.Printf("editor set to '%s'\n", name)
	return nil
}

// commandBase handles setting the base template or callsback the edit
// command to edit the base template based on the user flags
func (ctx *Context) commandBase(args []string) error {
	baseName := ctx.cfg.Base
	if len(args) == 0 && baseName == "" {
		return fmt.Errorf("no base template set")
	} else if len(args) == 0 {
		fmt.Printf("Your current base file is template: %v\n", baseName)
		return nil
	}

	name := args[0]
	if name == "" {
		return fmt.Errorf("no template name provided")
	}
	templ, err := config.FindTemplateByName(ctx.cfg, name)
	if err != nil {
		return fmt.Errorf("no template with that name exists")
	}
	ctx.cfg.Base = templ.Name
	if err := config.SaveConfig(ctx.cfg); err != nil {
		return err
	}
	fmt.Printf("base template set to '%s'\n", name)
	return nil
}

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

// commandAlias handles listing the available aliases
func (ctx *Context) commandAlias(args []string) error {
	fmt.Println("the available aliases are")
	for key, value := range aliasMap {
		fmt.Printf("%s -> %s\n", key, value)
	}
	return nil
}

// Helper function to resolve command aliases
func resolveCommand(name string) string {
	if primaryName, exists := aliasMap[name]; exists {
		return primaryName
	}
	return name
}

func openTemplateInEditor(editor, templPath string) error {
	cmd := exec.Command(editor, templPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
