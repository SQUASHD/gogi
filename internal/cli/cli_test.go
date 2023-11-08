package cli

import (
	"testing"
)

func TestSanitizeArgs(t *testing.T) {
	args := []string{"Gogi", "bASe", "TEST"}
	sanitizedArgs := sanitizeArgs(args)
	if len(sanitizedArgs) != 3 {
		t.Errorf("expected 2 args, got %d", len(sanitizedArgs))
	}
	if sanitizedArgs[0] != "gogi" {
		t.Errorf("expected gogi, got %s", sanitizedArgs[0])
	}
	if sanitizedArgs[1] != "base" {
		t.Errorf("expected base, got %s", sanitizedArgs[1])
	}
	if sanitizedArgs[2] != "test" {
		t.Errorf("expected test, got %s", sanitizedArgs[2])
	}
}
