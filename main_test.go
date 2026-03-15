package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ansible "github.com/arillso/go.ansible/v2"
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
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

func TestNormalizeSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "single value",
			input:    []string{"inv.yml"},
			expected: []string{"inv.yml"},
		},
		{
			name:     "comma-separated already split by cli",
			input:    []string{"inv1.yml", "inv2.yml"},
			expected: []string{"inv1.yml", "inv2.yml"},
		},
		{
			name:     "newline-separated from multiline YAML",
			input:    []string{"inv1.yml\ninv2.yml"},
			expected: []string{"inv1.yml", "inv2.yml"},
		},
		{
			name:     "newline with trailing newline",
			input:    []string{"inv1.yml\ninv2.yml\n"},
			expected: []string{"inv1.yml", "inv2.yml"},
		},
		{
			name:     "mixed whitespace and empty lines",
			input:    []string{"  inv1.yml \n\n  inv2.yml  \n"},
			expected: []string{"inv1.yml", "inv2.yml"},
		},
		{
			name:     "CRLF line endings",
			input:    []string{"inv1.yml\r\ninv2.yml\r\n"},
			expected: []string{"inv1.yml", "inv2.yml"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSlice(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d elements, got %d: %v", len(tt.expected), len(result), result)
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("element %d: expected %q, got %q", i, tt.expected[i], v)
				}
			}
		})
	}
}

func TestValidateParameters_MultilineInventory(t *testing.T) {
	tmpDir := t.TempDir()
	pbPath := createTempFile(t, tmpDir, "playbook.yml", "---\n- hosts: all\n")
	inv1 := createTempFile(t, tmpDir, "inventory1.yml", "all:\n  hosts:\n    localhost:\n")
	inv2 := createTempFile(t, tmpDir, "inventory2.yml", "all:\n  hosts:\n    localhost:\n")

	cmd := newTestCommand(func(ctx context.Context, c *cli.Command) error {
		return validateParameters(
			normalizeSlice(c.StringSlice("inventory")),
			normalizeSlice(c.StringSlice("playbook")),
			c.String("galaxy-file"),
		)
	})

	// Simulate multiline YAML input: newline-separated value in a single string.
	err := cmd.Run(context.Background(), []string{
		"test",
		"--playbook", pbPath,
		"--inventory", inv1 + "\n" + inv2,
	})
	if err != nil {
		t.Errorf("expected no error for newline-separated inventories, got: %v", err)
	}
}

func TestStartSSHAgent_ValidKey(t *testing.T) {
	if _, err := exec.LookPath("ssh-agent"); err != nil {
		t.Skip("ssh-agent not available")
	}
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not available")
	}

	// Generate a test RSA key.
	keyFile := filepath.Join(t.TempDir(), "test_key")
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", keyFile, "-N", "", "-q")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	keyContent, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	agent, err := startSSHAgent(string(keyContent), "")
	if err != nil {
		t.Fatalf("startSSHAgent failed: %v", err)
	}
	defer agent.stop()

	if agent.sock == "" {
		t.Error("expected SSH_AUTH_SOCK to be set")
	}
	if agent.pid == "" {
		t.Error("expected SSH_AGENT_PID to be set")
	}

	// Verify the key was added by listing keys.
	listCmd := exec.Command("ssh-add", "-l")
	listCmd.Env = append(os.Environ(), "SSH_AUTH_SOCK="+agent.sock)
	out, err := listCmd.Output()
	if err != nil {
		t.Fatalf("ssh-add -l failed: %v", err)
	}
	if !strings.Contains(string(out), "2048") {
		t.Errorf("expected key to be listed in agent, got: %s", string(out))
	}
}

func TestStartSSHAgent_InvalidKey(t *testing.T) {
	if _, err := exec.LookPath("ssh-agent"); err != nil {
		t.Skip("ssh-agent not available")
	}

	_, err := startSSHAgent("not-a-valid-key", "")
	if err == nil {
		t.Error("expected error for invalid key, got nil")
	}
}

func TestStartSSHAgent_CRLFKey(t *testing.T) {
	if _, err := exec.LookPath("ssh-agent"); err != nil {
		t.Skip("ssh-agent not available")
	}
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not available")
	}

	// Generate a test key and convert to CRLF.
	keyFile := filepath.Join(t.TempDir(), "test_key")
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", keyFile, "-N", "", "-q")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	keyContent, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	// Convert LF to CRLF.
	crlfKey := strings.ReplaceAll(string(keyContent), "\n", "\r\n")

	agent, err := startSSHAgent(crlfKey, "")
	if err != nil {
		t.Fatalf("startSSHAgent with CRLF key failed: %v", err)
	}
	defer agent.stop()

	if agent.sock == "" {
		t.Error("expected SSH_AUTH_SOCK to be set")
	}
}

func TestStartSSHAgent_PassphraseKey(t *testing.T) {
	if _, err := exec.LookPath("ssh-agent"); err != nil {
		t.Skip("ssh-agent not available")
	}
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not available")
	}

	// Generate a passphrase-protected RSA key.
	passphrase := "test-passphrase-123"
	keyFile := filepath.Join(t.TempDir(), "test_key")
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", keyFile, "-N", passphrase, "-q")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	keyContent, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	agent, err := startSSHAgent(string(keyContent), passphrase)
	if err != nil {
		t.Fatalf("startSSHAgent with passphrase failed: %v", err)
	}
	defer agent.stop()

	if agent.sock == "" {
		t.Error("expected SSH_AUTH_SOCK to be set")
	}

	// Verify the key was added by listing keys.
	listCmd := exec.Command("ssh-add", "-l")
	listCmd.Env = append(os.Environ(), "SSH_AUTH_SOCK="+agent.sock)
	out, err := listCmd.Output()
	if err != nil {
		t.Fatalf("ssh-add -l failed: %v", err)
	}
	if !strings.Contains(string(out), "2048") {
		t.Errorf("expected key to be listed in agent, got: %s", string(out))
	}
}

func TestStartSSHAgent_WrongPassphrase(t *testing.T) {
	if _, err := exec.LookPath("ssh-agent"); err != nil {
		t.Skip("ssh-agent not available")
	}
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not available")
	}

	// Generate a passphrase-protected key.
	keyFile := filepath.Join(t.TempDir(), "test_key")
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", keyFile, "-N", "correct-passphrase", "-q")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	keyContent, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	_, err = startSSHAgent(string(keyContent), "wrong-passphrase")
	if err == nil {
		t.Error("expected error for wrong passphrase, got nil")
	}
}

func TestSetupKnownHosts(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	content := "github.com ssh-rsa AAAAB3...\ngitlab.com ssh-ed25519 AAAAC3..."
	if err := setupKnownHosts(content); err != nil {
		t.Fatalf("setupKnownHosts failed: %v", err)
	}

	khPath := filepath.Join(tmpDir, ".ssh", "known_hosts")
	data, err := os.ReadFile(khPath)
	if err != nil {
		t.Fatalf("failed to read known_hosts: %v", err)
	}
	if !strings.Contains(string(data), "github.com") {
		t.Error("expected github.com in known_hosts")
	}
	if !strings.Contains(string(data), "gitlab.com") {
		t.Error("expected gitlab.com in known_hosts")
	}

	info, err := os.Stat(khPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %o", info.Mode().Perm())
	}

	dirInfo, err := os.Stat(filepath.Join(tmpDir, ".ssh"))
	if err != nil {
		t.Fatal(err)
	}
	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("expected 0700 on .ssh dir, got %o", dirInfo.Mode().Perm())
	}
}

func TestSetupKnownHosts_CRLF(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	content := "host1 ssh-rsa AAA\r\nhost2 ssh-rsa BBB\r\n"
	if err := setupKnownHosts(content); err != nil {
		t.Fatalf("setupKnownHosts with CRLF failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmpDir, ".ssh", "known_hosts"))
	if strings.Contains(string(data), "\r") {
		t.Error("expected CRLF to be normalized to LF")
	}
}

func TestWriteActionOutputs_Success(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "output")
	t.Setenv("GITHUB_OUTPUT", tmpFile)

	writeActionOutputs(nil)

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "status=success") {
		t.Errorf("expected status=success, got: %s", content)
	}
	if !strings.Contains(content, "exit_code=0") {
		t.Errorf("expected exit_code=0, got: %s", content)
	}
}

func TestWriteActionOutputs_Failed(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "output")
	t.Setenv("GITHUB_OUTPUT", tmpFile)

	writeActionOutputs(fmt.Errorf("some error"))

	data, _ := os.ReadFile(tmpFile)
	content := string(data)
	if !strings.Contains(content, "status=failed") {
		t.Errorf("expected status=failed, got: %s", content)
	}
	if !strings.Contains(content, "exit_code=1") {
		t.Errorf("expected exit_code=1, got: %s", content)
	}
}

func TestWriteActionOutputs_AnsibleError(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "output")
	t.Setenv("GITHUB_OUTPUT", tmpFile)

	writeActionOutputs(&ansible.AnsibleError{ExitCode: 2})

	data, _ := os.ReadFile(tmpFile)
	content := string(data)
	if !strings.Contains(content, "status=failed") {
		t.Errorf("expected status=failed, got: %s", content)
	}
	if !strings.Contains(content, "exit_code=2") {
		t.Errorf("expected exit_code=2, got: %s", content)
	}
}

func TestWriteActionOutputs_NoEnvVar(t *testing.T) {
	t.Setenv("GITHUB_OUTPUT", "")
	writeActionOutputs(nil)
}

func TestWriteStepSummary_Success(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "summary.md")
	t.Setenv("GITHUB_STEP_SUMMARY", tmpFile)

	writeStepSummary([]string{"site.yml", "deploy.yml"}, nil, 2*time.Minute+35*time.Second)

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read summary: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "site.yml") {
		t.Error("expected playbook name in summary")
	}
	if !strings.Contains(content, "Success") {
		t.Error("expected success status in summary")
	}
	if !strings.Contains(content, "2m 35s") {
		t.Errorf("expected '2m 35s' in summary, got: %s", content)
	}
}

func TestWriteStepSummary_NoEnvVar(t *testing.T) {
	t.Setenv("GITHUB_STEP_SUMMARY", "")
	// Should not panic or error.
	writeStepSummary([]string{"test.yml"}, nil, time.Second)
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m 30s"},
		{5*time.Minute + 10*time.Second, "5m 10s"},
	}
	for _, tt := range tests {
		if got := formatDuration(tt.d); got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestSSHAgentStop_Nil(t *testing.T) {
	// Calling stop on nil should not panic.
	var agent *sshAgent
	agent.stop()
}
