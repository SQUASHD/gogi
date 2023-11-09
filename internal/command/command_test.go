package command

import (
	"github.com/SQUASHD/gogi/internal/tests"
	"path/filepath"
	"testing"

	"github.com/SQUASHD/gogi/internal/structs"
)

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

func newTestContext(t *testing.T) (*Context, func()) {
	t.Helper()
	dir, cleanupFunc := tests.CreateTempDir(t)
	configPath := filepath.Join(dir, "config.json")
	cfg := createTestConfig(t, dir)
	ctx, _ := NewCommandContext(&cfg, dir, dir, configPath)
	return ctx, cleanupFunc
}

func TestQuickGogiWithValidBase(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	err := ctx.HandleQuickGogi()
	if err != nil {
		t.Errorf("Expected test to pass with valid base but got err: %v", err)
	}
}

func TestQuickGogiWithInvalidBase(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	ctx.cfg.Base = ""
	err := ctx.HandleQuickGogi()
	if err == nil {
		t.Errorf("Expected test to fail with invalid base but got err: %v", err)
	}
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

func TestRenameWithValidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandRename([]string{"test1", "test3"})
	if err != nil {
		t.Errorf("expected no error when calling rename with valid template name, but got err: %v", err)
	}
}

func TestRenameWithInvalidTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandRename([]string{"wrong", "test3"})
	if err == nil {
		t.Errorf("should not accept improper template name")
	}
}

func TestRenameToExistingTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandRename([]string{"test1", "test2"})
	if err == nil {
		t.Errorf("should not accept existing template name")
	}
}

func TestRenameWithMissingArgs(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandRename([]string{"test1"})
	if err == nil {
		t.Errorf("should not accept missing args")
	}
}

func TestRenameUpdatesBase(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()
	err := ctx.commandRename([]string{"test1", "test3"})
	if err != nil {
		t.Errorf("expected no error when calling rename with valid template name, but got err: %v", err)
	}
	if ctx.cfg.Base != "test3" {
		t.Errorf("expected base to be test3 but got %s", ctx.cfg.Base)
	}
}

func TestRenameUpdatesReferenceTemplate(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	err := ctx.commandRename([]string{"test1", "test3"})
	if err != nil {
		t.Fatalf("expected no error when calling rename with valid template name, but got: %v", err)
	}

	expectedPath := filepath.Join(ctx.projectDir, "test3.gitignore")

	// Cast to string because the path is of type string`json:"path"` of Template Type
	var resultPath string
	resultPath = ctx.cfg.Templates[0].Path

	if resultPath != expectedPath {
		t.Errorf("expected reference template to be %s but got %s", expectedPath, ctx.cfg.Templates[0].Path)
	}
}
