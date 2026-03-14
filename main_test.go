package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	cli "github.com/urfave/cli/v3"
)

// newTestCommand creates a CLI command with the same flags as the real app for testing.
func newTestCommand(action cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "execution-timeout",
				Value: 30,
			},
			&cli.StringFlag{
				Name: "galaxy-file",
			},
			&cli.BoolFlag{
				Name: "galaxy-force",
			},
			&cli.StringFlag{
				Name: "galaxy-api-key",
			},
			&cli.StringFlag{
				Name: "galaxy-api-server-url",
			},
			&cli.StringFlag{
				Name: "galaxy-collections-path",
			},
			&cli.BoolFlag{
				Name: "galaxy-disable-gpg-verify",
			},
			&cli.BoolFlag{
				Name: "galaxy-force-with-deps",
			},
			&cli.BoolFlag{
				Name: "galaxy-ignore-certs",
			},
			&cli.StringSliceFlag{
				Name: "galaxy-ignore-signature-status-codes",
			},
			&cli.StringFlag{
				Name: "galaxy-keyring",
			},
			&cli.BoolFlag{
				Name: "galaxy-offline",
			},
			&cli.BoolFlag{
				Name: "galaxy-pre",
			},
			&cli.IntFlag{
				Name: "galaxy-required-valid-signature-count",
			},
			&cli.StringFlag{
				Name: "galaxy-requirements-file",
			},
			&cli.StringFlag{
				Name: "galaxy-signature",
			},
			&cli.IntFlag{
				Name: "galaxy-timeout",
			},
			&cli.BoolFlag{
				Name: "galaxy-upgrade",
			},
			&cli.BoolFlag{
				Name: "galaxy-no-deps",
			},
			&cli.StringSliceFlag{
				Name: "inventory",
			},
			&cli.StringSliceFlag{
				Name: "playbook",
			},
			&cli.StringFlag{
				Name: "limit",
			},
			&cli.StringFlag{
				Name: "skip-tags",
			},
			&cli.StringFlag{
				Name: "start-at-task",
			},
			&cli.StringFlag{
				Name: "tags",
			},
			&cli.StringSliceFlag{
				Name: "extra-vars",
			},
			&cli.StringSliceFlag{
				Name: "module-path",
			},
			&cli.BoolFlag{
				Name: "check",
			},
			&cli.BoolFlag{
				Name: "diff",
			},
			&cli.BoolFlag{
				Name: "flush-cache",
			},
			&cli.BoolFlag{
				Name: "force-handlers",
			},
			&cli.BoolFlag{
				Name: "list-hosts",
			},
			&cli.BoolFlag{
				Name: "list-tags",
			},
			&cli.BoolFlag{
				Name: "list-tasks",
			},
			&cli.BoolFlag{
				Name: "syntax-check",
			},
			&cli.IntFlag{
				Name:  "forks",
				Value: 5,
			},
			&cli.StringFlag{
				Name: "vault-id",
			},
			&cli.StringFlag{
				Name: "vault-password",
			},
			&cli.IntFlag{
				Name: "verbose",
			},
			&cli.StringFlag{
				Name: "private-key",
			},
			&cli.StringFlag{
				Name: "private-key-file",
			},
			&cli.StringFlag{
				Name: "user",
			},
			&cli.StringFlag{
				Name: "connection",
			},
			&cli.IntFlag{
				Name: "timeout",
			},
			&cli.StringFlag{
				Name: "ssh-common-args",
			},
			&cli.StringFlag{
				Name: "sftp-extra-args",
			},
			&cli.StringFlag{
				Name: "scp-extra-args",
			},
			&cli.StringFlag{
				Name: "ssh-extra-args",
			},
			&cli.BoolFlag{
				Name: "become",
			},
			&cli.StringFlag{
				Name: "become-method",
			},
			&cli.StringFlag{
				Name: "become-user",
			},
			&cli.BoolFlag{
				Name: "ask-become-pass",
			},
			&cli.BoolFlag{
				Name: "ask-pass",
			},
			&cli.BoolFlag{
				Name: "step",
			},
			&cli.StringFlag{
				Name: "ssh-transfer-method",
			},
			&cli.StringFlag{
				Name: "output-callback",
			},
			&cli.BoolFlag{
				Name: "no-color",
			},
			&cli.StringFlag{
				Name: "vault-password-file",
			},
			&cli.BoolFlag{
				Name: "ask-vault-pass",
			},
			&cli.StringFlag{
				Name: "fact-path",
			},
			&cli.StringFlag{
				Name: "fact-caching",
			},
			&cli.IntFlag{
				Name: "fact-caching-timeout",
			},
			&cli.StringFlag{
				Name: "callbacks-enabled",
			},
			&cli.IntFlag{
				Name: "poll-interval",
			},
			&cli.StringFlag{
				Name: "gather-subset",
			},
			&cli.IntFlag{
				Name: "gather-timeout",
			},
			&cli.StringFlag{
				Name: "strategy-plugin",
			},
			&cli.IntFlag{
				Name: "max-fail-percentage",
			},
			&cli.BoolFlag{
				Name: "any-errors-fatal",
			},
			&cli.StringFlag{
				Name: "config-file",
			},
			&cli.StringFlag{
				Name: "temp-dir",
			},
		},
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

	// Verify error wraps ErrInvalidParameter
	if !errorContains(err, "invalid parameter provided") {
		t.Errorf("expected error to contain ErrInvalidParameter message, got: %v", err)
	}
	if !errorContains(err, "playbook file does not exist") {
		t.Errorf("expected error to mention playbook file, got: %v", err)
	}
}

// errorContains checks if an error message contains the given substring.
func errorContains(err error, substr string) bool {
	if err == nil {
		return false
	}
	return contains(err.Error(), substr)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
