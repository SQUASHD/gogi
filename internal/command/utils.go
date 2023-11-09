package command

import (
	"os"
	"os/exec"
)

// Helper function to resolve command aliases
func resolveCommand(name string) string {
	if primaryName, exists := aliasMap[name]; exists {
		return primaryName
	}
	return name
}

// openTemplateInEditor opens the template in the user's editor
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
