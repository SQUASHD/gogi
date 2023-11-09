package command

import (
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
	"r": "rename",
}

// Context holds the state and provides methods to execute CLI commands
type Context struct {
	cfg        *structs.TemplateConfig
	cwd        string
	commands   map[string]cliCommand
	projectDir string
	configPath string
}

// cliCommand represents a command in the CLI
type cliCommand struct {
	name        string
	description string
	helpExample string
	callback    func(*Context, []string) error
}

// NewCommandContext initializes a new command context with the given configuration
// and current working directory
func NewCommandContext(cfg *structs.TemplateConfig, cwd, projectDir, configPath string) (*Context, error) {
	ctx := &Context{
		cfg:        cfg,
		cwd:        cwd,
		projectDir: projectDir,
		configPath: configPath,
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
			helpExample: "gogi create template-name",
			callback:    (*Context).commandCreate,
		},
		"delete": {
			name:        "delete",
			description: "Delete an existing gitignore alias",
			helpExample: "gogi delete template-name",
			callback:    (*Context).commandDelete,
		},
		"list": {
			name:        "list",
			description: "List all the templates",
			helpExample: "gogi list",
			callback:    (*Context).commandList,
		},
		"generate": {
			name:        "generate",
			description: "Generate a gitignore file from the given template",
			helpExample: "gogi generate template-name",
			callback:    (*Context).commandGenerate,
		},
		"edit": {
			name:        "edit",
			description: "Edit an existing template",
			helpExample: "gogi edit template-name",
			callback:    (*Context).commandEdit,
		},
		"append": {
			name:        "append",
			description: "Append a template to an existing gitignore file",
			helpExample: "gogi append template-name",
			callback:    (*Context).commandAppend,
		},
		"help": {
			name:        "help",
			description: "Display help message, or help for a specific command",
			helpExample: "gogi help [command]",
			callback:    (*Context).commandHelp,
		},
		"editor": {
			name:        "editor",
			description: "Set the editor to use for editing templates",
			helpExample: "gogi editor editor-name",
			callback:    (*Context).commandEditor,
		},
		"base": {
			name:        "base",
			description: "set the base template that you call with gogi with no args",
			helpExample: "gogi base template-name",
			callback:    (*Context).commandBase,
		},
		"alias": {
			name:        "alias",
			description: "show the list of avaiable command aliases",
			helpExample: "gogi alias",
			callback:    (*Context).commandAlias,
		},
		"rename": {
			name:        "rename",
			description: "rename a template",
			helpExample: "gogi rename old-name new-name",
			callback:    (*Context).commandRename,
		},
	}
}
