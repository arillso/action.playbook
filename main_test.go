package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/urfave/cli/v3"
)

// newTestCommand creates a CLI command reusing appFlags from main.go for testing.
func newTestCommand(action cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   "test",
		Flags:  appFlags,
		Action: action,
	}
}

// createTempFile creates a temporary file with the given content and returns its path.
func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to create temp file %s: %v", name, err)
	}
	return path
}

func TestValidateParameters_ValidFiles(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", invPath,
	})
	if err != nil {
		t.Errorf("expected no error for valid files, got: %v", err)
	}
}

func TestValidateParameters_MissingPlaybook(t *testing.T) {
	tmpDir := t.TempDir()
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", filepath.Join(tmpDir, "nonexistent.yml"),
		"--inventory", invPath,
	})
	if err == nil {
		t.Error("expected error for missing playbook, got nil")
	}
}

func TestValidateParameters_MissingInventory(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", filepath.Join(tmpDir, "nonexistent.yml"),
	})
	if err == nil {
		t.Error("expected error for missing inventory, got nil")
	}
}

func TestValidateParameters_MissingGalaxyFile(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", invPath,
		"--galaxy-file", filepath.Join(tmpDir, "nonexistent_requirements.yml"),
	})
	if err == nil {
		t.Error("expected error for missing galaxy file, got nil")
	}
}

func TestValidateParameters_ValidGalaxyFile(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")
	galaxyPath := createTempFile(t, tmpDir, "requirements.yml", "---\nroles: []\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", invPath,
		"--galaxy-file", galaxyPath,
	})
	if err != nil {
		t.Errorf("expected no error for valid galaxy file, got: %v", err)
	}
}

func TestValidateParameters_NoGalaxyFile(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", invPath,
	})
	if err != nil {
		t.Errorf("expected no error when galaxy file is not specified, got: %v", err)
	}
}

func TestValidateParameters_MultiplePlaybooks(t *testing.T) {
	tmpDir := t.TempDir()
	pb1 := createTempFile(t, tmpDir, "playbook1.yml", "---\n- hosts: all\n")
	pb2 := createTempFile(t, tmpDir, "playbook2.yml", "---\n- hosts: all\n")
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pb1,
		"--playbook", pb2,
		"--inventory", invPath,
	})
	if err != nil {
		t.Errorf("expected no error for multiple valid playbooks, got: %v", err)
	}
}

func TestValidateParameters_MultiplePlaybooksOneMissing(t *testing.T) {
	tmpDir := t.TempDir()
	pb1 := createTempFile(t, tmpDir, "playbook1.yml", "---\n- hosts: all\n")
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pb1,
		"--playbook", filepath.Join(tmpDir, "missing.yml"),
		"--inventory", invPath,
	})
	if err == nil {
		t.Error("expected error when one of multiple playbooks is missing, got nil")
	}
}

func TestValidateParameters_MultipleInventories(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")
	inv1 := createTempFile(t, tmpDir, "inventory1.yml", "all:\n  hosts:\n    localhost:\n")
	inv2 := createTempFile(t, tmpDir, "inventory2.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", inv1,
		"--inventory", inv2,
	})
	if err != nil {
		t.Errorf("expected no error for multiple valid inventories, got: %v", err)
	}
}

func TestValidateParameters_ErrorWrapping(t *testing.T) {
	tmpDir := t.TempDir()
	invPath := createTempFile(t, tmpDir, "inventory.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(c)
	})

	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", filepath.Join(tmpDir, "missing.yml"),
		"--inventory", invPath,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid parameter provided") {
		t.Errorf("expected error to contain ErrInvalidParameter message, got: %v", err)
	}
	if !strings.Contains(err.Error(), "playbook file does not exist") {
		t.Errorf("expected error to mention playbook file, got: %v", err)
	}
}
