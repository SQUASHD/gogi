package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SQUASHD/gogi/internal/structs"
)

func createTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	_, err = os.Create(filepath.Join(tempDir, "test1.gitignore"))
	if err != nil {
		t.Fatalf("Failed to create test1.gitignore file: %v", err)
	}
	_, err = os.Create(filepath.Join(tempDir, "test2.gitignore"))
	if err != nil {
		t.Fatalf("Failed to create test2.gitignore file: %v", err)
	}

	return tempDir, func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("Failed to remove temp directory: %v", err)
		}
	}
}

func createTestConfig(t *testing.T, folderPath string) structs.TemplateConfig {
	t.Helper()
	testConfig := structs.TemplateConfig{
		Editor: "nvim",
		Base:   "test1",
		Templates: []structs.Template{
			{
				Name: "test1",
				Path: filepath.Join(folderPath, "test1.gitignore"),
			},
			{
				Name: "test2",
				Path: filepath.Join(folderPath, "test2.gitignore"),
			},
		},
	}
	return testConfig
}

func newTestContext(t *testing.T) (*CommandContext, func()) {
	t.Helper()
	testProjectDir, cleanupFunc := createTempDir(t)

	testConfig := createTestConfig(t, testProjectDir)

	ctx := &CommandContext{
		cfg:        &testConfig,
		cwd:        testProjectDir,
		projectDir: testProjectDir,
	}
	ctx.commands = ctx.getCommands()

	return ctx, cleanupFunc
}

func TestDeleteCommandWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	err := ctx.commandDelete([]string{"doesnotexist"})
	if err == nil {
		t.Errorf("should not accept improper template name")
	}
}

func TestDeleteCommandWithValidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	err := ctx.commandDelete([]string{"test2"})
	if err != nil {
		t.Errorf("Expected test to pass with valid template name but got err: %v", err)
	}

	if ctx.cfg.Base != "test1" {
		t.Errorf("Expected base to be test1 but got %s", ctx.cfg.Base)
	}

	if len(ctx.cfg.Templates) != 1 {
		t.Errorf("Expected templates to have length 1 but got %d", len(ctx.cfg.Templates))
	}
}

func TestDeleteCommandWithBaseTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	err := ctx.commandDelete([]string{"test1"})
	if err != nil {
		t.Errorf("Expected test to pass with valid template name but got err: %v", err)
	}

	if ctx.cfg.Base != "" {
		t.Errorf("Expected base to be empty but got %s", ctx.cfg.Base)
	}
	if len(ctx.cfg.Templates) != 1 {
		t.Errorf("Expected templates to have length 1 but got %d", len(ctx.cfg.Templates))
	}
}

func TestBaseCommandWithSameTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandBase([]string{"test1"})
	if err != nil {
		t.Errorf("expected no error when setting base to same template")
	}
}

func TestBaseCommandWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandBase([]string{"test3"})
	if err == nil {
		t.Errorf("expected to test with valid template name, but got err: %v", err)
	}
}

func TestBaseCommandWithValidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandBase([]string{"test2"})
	if err != nil {
		t.Errorf("expected to test with valid template name, but got err: %v", err)
	}
}

func TestCreateCommandWithExistingTemplateName(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandCreate([]string{"test1"})
	if err == nil {
		t.Errorf("should not accept existing template name")
	}
}

func TestGenerateCommandWithValidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandGenerate([]string{"test1"})
	if err != nil {
		t.Errorf("expected to test with valid template name, but got err: %v", err)
	}
}

func TestGenerateCommandWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandGenerate([]string{"wrong"})
	if err == nil {
		t.Errorf("should not accept improper template name")
	}
}

func TestHelpWithNoArgs(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandHelp([]string{})
	if err != nil {
		t.Errorf("expected no error when calling help with no args, but got err: %v", err)
	}
}

func TestHelpWithValidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandHelp([]string{"h"})
	if err != nil {
		t.Errorf("expected no error when calling help with valid template name, but got err: %v", err)
	}
}

func TestHelpWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandHelp([]string{"wrong"})
	if err == nil {
		t.Errorf("should not accept improper template name")
	}
}

func TestCommandEditWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandEdit([]string{"wrong"})
	if err == nil {
		t.Errorf("should not accept improper template name")
	}
}

func TestCommandEditorWithNoArgs(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandEditor([]string{})
	if err == nil {
		t.Errorf("expected no error when calling editor with no args, but got err: %v", err)
	}
}

func TestCommandEditorWithValidEditor(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandEditor([]string{"nvim"})
	if err != nil {
		t.Errorf("expected no error when calling editor with valid template name, but got err: %v", err)
	}
}

func TestCommandListWithNoArgs(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandList([]string{})
	if err != nil {
		t.Errorf("expected no error when calling list with no args, but got err: %v", err)
	}
}

func TestCommandBaseWithBaseTemplateAndNoArgs(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandBase([]string{})
	if err != nil {
		t.Errorf("expected no error when calling base with no args, but got err: %v", err)
	}
}

func TestCommandBaseWithNoBaseTemplateAndNoArgs(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	ctx.cfg.Base = ""
	err := ctx.commandBase([]string{})
	if err == nil {
		t.Errorf("expected error when calling base with no base template")
	}
}

func TestCommandCreateThenDelete(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandCreate([]string{"test3"})
	if err != nil {
		t.Errorf("expected no error when calling add with valid template name, but got err: %v", err)
	}
	err = ctx.commandDelete([]string{"test3"})
	if err != nil {
		t.Errorf("expected no error when calling delete with valid template name, but got err: %v", err)
	}
}

func TestCommandCreateThenGenerate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandCreate([]string{"test3"})
	if err != nil {
		t.Errorf("expected no error when calling add with valid template name, but got err: %v", err)
	}
	err = ctx.commandGenerate([]string{"test3"})
	if err != nil {
		t.Errorf("expected no error when calling generate with valid template name, but got err: %v", err)
	}
}

func TestCommandCreateThenBase(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandCreate([]string{"test3"})
	if err != nil {
		t.Errorf("expected no error when calling add with valid template name, but got err: %v", err)
	}
	err = ctx.commandBase([]string{"test3"})
	if err != nil {
		t.Errorf("expected no error when calling base with valid template name, but got err: %v", err)
	}
}

func TestGenerateTwice(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandGenerate([]string{"test1"})
	if err != nil {
		t.Errorf("expected no error when calling generate with valid template name, but got err: %v", err)
	}
	err = ctx.commandGenerate([]string{"test1"})
	if err == nil {
		t.Errorf("expected error when calling generate twice")
	}
}

func TestCommandAppendWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandGenerate([]string{"test1"})
	if err != nil {
		t.Errorf("expected no error when calling generate with valid template name, but got err: %v", err)
	}
	err = ctx.commandAppend([]string{"wrong"})
	if err == nil {
		t.Errorf("should not accept improper template name")
	}
}

func TestCommandAppendWithValidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandGenerate([]string{"test1"})
	if err != nil {
		t.Errorf("expected no error when calling generate with valid template name, but got err: %v", err)
	}
	err = ctx.commandAppend([]string{"test2"})
	if err != nil {
		t.Errorf("expected no error when calling append with valid template name, but got err: %v", err)
	}
}
