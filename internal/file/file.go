package file

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenTemplateInEditor(editor, templPath string) error {

	if _, err := os.Stat(templPath); os.IsNotExist(err) {
		return fmt.Errorf("no template found at: %s", templPath)
	}
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

func DeleteTemplate(templPath string) error {
	if err := os.Remove(templPath); err != nil {
		return err
	}
	return nil
}
