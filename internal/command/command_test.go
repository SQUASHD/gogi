package command

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
	file1, err := os.Create(filepath.Join(tempDir, "test1.gitignore"))
	if err != nil {
		t.Fatalf("Failed to create test1.gitignore file: %v", err)
	}
	file1.Close()

	file2, err := os.Create(filepath.Join(tempDir, "test2.gitignore"))
	if err != nil {
		t.Fatalf("Failed to create test2.gitignore file: %v", err)
	}
	file2.Close()

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
		Editor: "nano",
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
	dir, cleanupFunc := createTempDir(t)
	configPath := filepath.Join(dir, "config.json")
	cfg := createTestConfig(t, dir)
	ctx, _ := NewCommandContext(&cfg, dir, dir, configPath)
	return ctx, cleanupFunc
}

func addGitIgnoreToTestDir(t *testing.T, ctx *Context) {
	t.Helper()
	gitignorePath := filepath.Join(ctx.cwd, ".gitignore")
	destinationFile, err := os.Create(gitignorePath)
	if err != nil {
		t.Fatalf("unable to create .gitignore file: %v", err)
	}
	defer destinationFile.Close()
}

func TestAppendCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		giExist bool
		wantErr bool
	}{
		{"append valid", []string{"test2"}, true, false},
		{"append invalid", []string{"invalid"}, true, true},
		{"append valid no gitignore", []string{"test2"}, false, true},
		{"append invalid no gitignore", []string{"invalid"}, false, true},
		{"append no args", []string{}, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			if tt.giExist {
				addGitIgnoreToTestDir(t, ctx)
			}

			err := ctx.commandAppend(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandAppend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		expected string
	}{
		{"set valid base", []string{"test2"}, false, "test2"},
		{"set invalid base", []string{"invalid"}, true, "test1"},
		{"set same base", []string{"test1"}, false, "test1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			err := ctx.commandBase(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandBase() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && ctx.cfg.Base != tt.expected {
				t.Errorf("Expected base to be %s but got %s", tt.expected, ctx.cfg.Base)
			}
		})
	}
}

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		expectedLen int
	}{
		{"create valid name", []string{"test3"}, false, 3},
		{"create existing name", []string{"test1"}, true, 2},
		{"create malformed arg", []string{""}, true, 2},
		{"create no args", []string{}, true, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			err := ctx.commandCreate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandDelete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(ctx.cfg.Templates) != tt.expectedLen {
				t.Errorf("Expected templates to have length %d but got %d", tt.expectedLen, len(ctx.cfg.Templates))
			}
		})
	}
}

func TestRenameCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		pathChanged bool
	}{
		{"no args", []string{}, true, false},
		{"missing args", []string{"test1"}, true, false},
		{"rename to existing template", []string{"test1", "test2"}, true, false},
		{"rename to new template", []string{"test1", "test3"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			originalPath := ctx.cfg.Templates[0].Path
			err := ctx.commandRename(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("comandRename() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.pathChanged && ctx.cfg.Templates[0].Path == originalPath {
				t.Errorf("Expected path to change but it did not")
			}
		})
	}
}

func TestDeleteCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantErr      bool
		expectedBase string
		expectedLen  int
	}{
		{"delete valid", []string{"test2", "--force"}, false, "test1", 1},
		{"delete invalid", []string{"invalid", "--force"}, true, "test1", 2},
		{"delete base", []string{"test1", "--force"}, false, "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			err := ctx.commandDelete(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandDelete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if ctx.cfg.Base != tt.expectedBase {
				t.Errorf("Expected base to be %s but got %s", tt.expectedBase, ctx.cfg.Base)
			}

			if len(ctx.cfg.Templates) != tt.expectedLen {
				t.Errorf("Expected templates to have length %d but got %d", tt.expectedLen, len(ctx.cfg.Templates))
			}
		})
	}
}

func TestEditCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		editorSet bool
		wantErr   bool
	}{
		{"edit no args with editor", []string{}, true, true},
		{"edit no args without editor", []string{}, false, true},
		{"edit malformed args with editor", []string{""}, true, true},
		{"edit malformed args without editor", []string{""}, false, true},

		// Some tests are diabled due to the command opening the editor
		//{"edit invalid template", []string{"invalid"}, true, false},
		//{"edit valid template", []string{"test1"}, true, false},
		{"edit valid template without editor", []string{"test1"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			if !tt.editorSet {
				ctx.cfg.Editor = ""
			}

			err := ctx.commandEdit(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandEditor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditorCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"editor no args", []string{}, false},
		{"editor malformed args", []string{""}, true},
		{"editor same editor", []string{"nano"}, false},
		{"edidtor new editor", []string{"code"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()

			err := ctx.commandEditor(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandEditor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandGenerate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		giExist bool
		wantErr bool
	}{
		{"generate valid", []string{"test2", "--force"}, false, false},
		{"generate invalid", []string{"invalid", "--force"}, false, true},
		{"generate valid with gitignore", []string{"test2", "--force"}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()
			if tt.giExist {
				addGitIgnoreToTestDir(t, ctx)
			}
			err := ctx.commandGenerate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"no args", []string{}, false},
		{"valid args", []string{"help"}, false},
		{"invalid args", []string{"invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()
			err := ctx.commandHelp(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandHelp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuickGogi(t *testing.T) {
	tests := []struct {
		name      string
		validBase bool
		wantErr   bool
	}{
		{"gogi with no base", true, false},
		{"gogi with base", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cleanup := newTestContext(t)
			defer cleanup()
			if !tt.validBase {
				ctx.cfg.Base = ""
			}

			err := ctx.HandleQuickGogi()
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleQuickGogi() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
