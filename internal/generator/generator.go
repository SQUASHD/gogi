package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GenerateGitignore creates or overwrites a .gitignore file at gitignorePath
// using the template found at templPath.
func GenerateGitignore(templPath, cwd string) error {
	sourceFile, err := os.Open(templPath)
	if err != nil {
		return fmt.Errorf("unable to open template file: %w", err)
	}
	defer sourceFile.Close()

	gitignorePath := filepath.Join(cwd, ".gitignore")
	destinationFile, err := os.Create(gitignorePath)
	if err != nil {
		return fmt.Errorf("unable to create .gitignore file: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("unable to write to .gitignore file: %w", err)
	}

	return nil
}

// AppendTemplate appends the contents of the template found at templPath
func AppendTemplate(cwd, templPath string) error {
	template, err := os.Open(templPath)
	if err != nil {
		return err
	}
	defer template.Close()

	giPath := filepath.Join(cwd, ".gitignore")
	file, err := os.OpenFile(giPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, template)
	if err != nil {
		return err
	}

	return nil
}

// CreateEmptyTemplateFile creates an empty template file at templPath
func CreateEmptyTemplateFile(projectDir, templName string) error {
	filename := templName + ".gitignore"
	templPath := filepath.Join(projectDir, filename)

	file, err := os.Create(templPath)
	if err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func DeleteTemplateFile(projectDir, templName string) error {
	fileName := templName + ".gitignore"
	err := os.Chdir(projectDir)
	if err != nil {
		return fmt.Errorf("error changing directory to delete file: %v", err)
	}
	if err = os.Remove(fileName); err != nil {
		return fmt.Errorf("error deleting template: %v", err)
	}
	return nil
}

// DoesGitignoreExist checks whether a .gitignore file exists
// in the current working directory.
func DoesGitignoreExist(currDir string) (bool, error) {
	gitignorePath := filepath.Join(currDir, ".gitignore")
	_, err := os.Stat(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking for .gitignore file: %v", err)
	}
	return true, nil
}

// CheckWhetherTemplateExists checks whether a template exists at the
// at the given file path
func CheckWhetherTemplateExists(templPath string) error {
	_, err := os.Stat(templPath)
	if err != nil {
		return fmt.Errorf("template does not exist at: %s", err)
	}
	return nil
}
