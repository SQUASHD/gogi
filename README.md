## Gogi: A simple .gitignore template manager

![the installation process for gogi](./gifs/gogi_init.gif)

### Installation
Deploy Gogi with the following command:
```bash
go install github.com/SQUASHD/gogi@latest
```

### Configuration
Initialize Gogi to set up a dedicated template directory within your `.config` folder:
```bash
gogi init
```

### Template Management
![using gogi](./gifs/gogi_editor.gif)
Craft a new, blank template 
```bash
gogi create <template-name> [-e open in editor] [-b set as base]
```

Specify your preferred editor for template customization:
```bash
gogi editor <editor-name>
```

Modify an existing template to suit your project needs:
```bash
gogi edit <template-name>
```

Establish a default template for new projects:
```bash
gogi base <template-name>
```

Delete an outdated template
```bash
gogi delete <template-name> [--force will override the are you sure prompt]
```

### Generate .gitignore
Generate a .gitignore file using your base template directly in your current project directory:
```bash
gogi
```

If you need a different template you can use gogi generate to generate a new 
.gitignore file

```bash
gogi generate <template-name> [--force will overwrite the current .gitignore]
```

Or append the current .gitignore with a different template

```bash
gogi append <template-name>
```

### Assistance



```txt
  create: Create a new template
  delete: Delete an existing gitignore alias
  editor: Set the editor to use for editing templates
   alias: show the list of avaiable command aliases
    help: Display help message, or help for a specific command
    base: set the base template that you call with gogi with no args
  rename: rename a template
    list: List all the templates
generate: Generate a gitignore file from the given template
    edit: Edit an existing template
  append: Append a template to an existing gitignore file
```

Most commands have an alias corresponding to their first letter
```
d -> delete
a -> append
r -> rename
h -> help
c -> create
l -> list
g -> generate
e -> edit
b -> base
```

See the whole suite of Gogi commands at any point
```bash
gogi help
```
