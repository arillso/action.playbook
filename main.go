// Package main implements the entry point for the Ansible Playbook Wrapper.
// This application executes Ansible playbooks and supports various configuration options.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	ansible "github.com/arillso/go.ansible/v2"
	"github.com/joho/godotenv"
	cli "github.com/urfave/cli/v3"
)

// Errors defined for better error handling
var (
	ErrPlaybookExecution = errors.New("playbook execution failed")
	ErrConfigLoad        = errors.New("failed to load configuration")
	ErrInvalidParameter  = errors.New("invalid parameter provided")
)

// appFlags defines all CLI flags for the application.
// Exported as a package-level variable so tests can reuse the same flag definitions.
var appFlags = []cli.Flag{
	&cli.IntFlag{
		Name:    "execution-timeout",
		Usage:   "Timeout in minutes for the playbook execution (default: 30)",
		Value:   30,
		Sources: cli.EnvVars("ANSIBLE_EXECUTION_TIMEOUT", "INPUT_EXECUTION_TIMEOUT", "PLUGIN_EXECUTION_TIMEOUT"),
	},
	// Galaxy-related options
	&cli.StringFlag{
		Name:    "galaxy-file",
		Usage:   "Path to the Ansible Galaxy requirements file",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_FILE", "INPUT_GALAXY_FILE", "PLUGIN_GALAXY_FILE"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-force",
		Usage:   "Force reinstallation of roles or collections from the Galaxy file",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_FORCE", "INPUT_GALAXY_FORCE", "PLUGIN_GALAXY_FORCE"),
	},
	&cli.StringFlag{
		Name:    "galaxy-api-key",
		Usage:   "API key for authenticating with Ansible Galaxy",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_API_KEY"),
	},
	&cli.StringFlag{
		Name:    "galaxy-api-server-url",
		Usage:   "URL of the Ansible Galaxy API server",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_API_SERVER_URL"),
	},
	&cli.StringFlag{
		Name:    "galaxy-collections-path",
		Usage:   "Path to the directory where Galaxy collections are stored",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_COLLECTIONS_PATH"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-disable-gpg-verify",
		Usage:   "Disable GPG signature verification for Galaxy operations",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_DISABLE_GPG_VERIFY"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-force-with-deps",
		Usage:   "Force installation of collections including their dependencies",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_FORCE_WITH_DEPS"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-ignore-certs",
		Usage:   "Ignore SSL certificate validation for Galaxy requests",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_IGNORE_CERTS"),
	},
	&cli.StringSliceFlag{
		Name:    "galaxy-ignore-signature-status-codes",
		Usage:   "Comma-separated list of HTTP status codes to ignore during signature validation",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_IGNORE_SIGNATURE_STATUS_CODES"),
	},
	&cli.StringFlag{
		Name:    "galaxy-keyring",
		Usage:   "Path to the GPG keyring file for Galaxy",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_KEYRING"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-offline",
		Usage:   "Enable offline mode to prevent requests to Ansible Galaxy",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_OFFLINE"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-pre",
		Usage:   "Allow installation of pre-release versions from Galaxy",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_PRE"),
	},
	&cli.IntFlag{
		Name:    "galaxy-required-valid-signature-count",
		Usage:   "Required number of valid GPG signatures for Galaxy content",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_REQUIRED_VALID_SIGNATURE_COUNT"),
	},
	&cli.StringFlag{
		Name:    "galaxy-requirements-file",
		Usage:   "Path to the Ansible Galaxy requirements file",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_REQUIREMENTS_FILE"),
	},
	&cli.StringFlag{
		Name:    "galaxy-signature",
		Usage:   "Specific GPG signature to verify for Galaxy content",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_SIGNATURE"),
	},
	&cli.IntFlag{
		Name:    "galaxy-timeout",
		Usage:   "Timeout (in seconds) for Galaxy operations",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_TIMEOUT"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-upgrade",
		Usage:   "Automatically upgrade Galaxy collections to the latest version",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_UPGRADE"),
	},
	&cli.BoolFlag{
		Name:    "galaxy-no-deps",
		Usage:   "Disable automatic dependency resolution for Galaxy",
		Sources: cli.EnvVars("ANSIBLE_GALAXY_NO_DEPS"),
	},
	// Inventory and playbook options
	&cli.StringSliceFlag{
		Name:     "inventory",
		Aliases:  []string{"i"},
		Usage:    "Path to one or more inventory files for Ansible",
		Sources:  cli.EnvVars("ANSIBLE_INVENTORY", "INPUT_INVENTORY", "PLUGIN_INVENTORY"),
		Required: true,
	},
	&cli.StringSliceFlag{
		Name:     "playbook",
		Aliases:  []string{"p"},
		Usage:    "List of playbooks to run",
		Sources:  cli.EnvVars("ANSIBLE_PLAYBOOK", "INPUT_PLAYBOOK", "PLUGIN_PLAYBOOK"),
		Required: true,
	},
	&cli.StringFlag{
		Name:    "limit",
		Aliases: []string{"l"},
		Usage:   "Limit playbook execution to a specific host group",
		Sources: cli.EnvVars("ANSIBLE_LIMIT", "INPUT_LIMIT", "PLUGIN_LIMIT"),
	},
	&cli.StringFlag{
		Name:    "skip-tags",
		Usage:   "Skip plays and tasks that match the given tags",
		Sources: cli.EnvVars("ANSIBLE_SKIP_TAGS", "INPUT_SKIP_TAGS", "PLUGIN_SKIP_TAGS"),
	},
	&cli.StringFlag{
		Name:    "start-at-task",
		Usage:   "Start playbook execution at the task with the given name",
		Sources: cli.EnvVars("ANSIBLE_START_AT_TASK", "INPUT_START_AT_TASK", "PLUGIN_START_AT_TASK"),
	},
	&cli.StringFlag{
		Name:    "tags",
		Aliases: []string{"t"},
		Usage:   "Run only tasks and plays with the specified tags",
		Sources: cli.EnvVars("ANSIBLE_TAGS", "INPUT_TAGS", "PLUGIN_TAGS"),
	},
	&cli.StringSliceFlag{
		Name:    "extra-vars",
		Aliases: []string{"e"},
		Usage:   "Set additional variables in key=value format",
		Sources: cli.EnvVars("ANSIBLE_EXTRA_VARS", "INPUT_EXTRA_VARS", "PLUGIN_EXTRA_VARS"),
	},
	&cli.StringSliceFlag{
		Name:    "module-path",
		Aliases: []string{"M"},
		Usage:   "Prepend directories to the module library path",
		Sources: cli.EnvVars("ANSIBLE_MODULE_PATH", "INPUT_MODULE_PATH", "PLUGIN_MODULE_PATH"),
	},
	&cli.BoolFlag{
		Name:    "check",
		Aliases: []string{"C"},
		Usage:   "Perform a dry run without making any changes",
		Sources: cli.EnvVars("ANSIBLE_CHECK", "INPUT_CHECK", "PLUGIN_CHECK"),
	},
	&cli.BoolFlag{
		Name:    "diff",
		Aliases: []string{"D"},
		Usage:   "Show the differences in files or templates when changes occur",
		Sources: cli.EnvVars("ANSIBLE_DIFF", "INPUT_DIFF", "PLUGIN_DIFF"),
	},
	&cli.BoolFlag{
		Name:    "flush-cache",
		Usage:   "Clear the fact cache for all hosts in the inventory",
		Sources: cli.EnvVars("ANSIBLE_FLUSH_CACHE", "INPUT_FLUSH_CACHE", "PLUGIN_FLUSH_CACHE"),
	},
	&cli.BoolFlag{
		Name:    "force-handlers",
		Usage:   "Run all handlers even if a task fails",
		Sources: cli.EnvVars("ANSIBLE_FORCE_HANDLERS", "INPUT_FORCE_HANDLERS", "PLUGIN_FORCE_HANDLERS"),
	},
	&cli.BoolFlag{
		Name:    "list-hosts",
		Usage:   "Display a list of matching hosts",
		Sources: cli.EnvVars("ANSIBLE_LIST_HOSTS", "INPUT_LIST_HOSTS", "PLUGIN_LIST_HOSTS"),
	},
	&cli.BoolFlag{
		Name:    "list-tags",
		Usage:   "List all available tags",
		Sources: cli.EnvVars("ANSIBLE_LIST_TAGS", "INPUT_LIST_TAGS", "PLUGIN_LIST_TAGS"),
	},
	&cli.BoolFlag{
		Name:    "list-tasks",
		Usage:   "List all tasks that would be executed",
		Sources: cli.EnvVars("ANSIBLE_LIST_TASKS", "INPUT_LIST_TASKS", "PLUGIN_LIST_TASKS"),
	},
	&cli.BoolFlag{
		Name:    "syntax-check",
		Usage:   "Perform a syntax check on the playbook without executing it",
		Sources: cli.EnvVars("ANSIBLE_SYNTAX_CHECK", "INPUT_SYNTAX_CHECK", "PLUGIN_SYNTAX_CHECK"),
	},
	&cli.IntFlag{
		Name:    "forks",
		Aliases: []string{"f"},
		Usage:   "Number of parallel processes to use during playbook execution",
		Value:   5,
		Sources: cli.EnvVars("ANSIBLE_FORKS", "INPUT_FORKS", "PLUGIN_FORKS"),
	},
	// Vault and authentication options
	&cli.StringFlag{
		Name:    "vault-id",
		Usage:   "Identity to use when accessing an Ansible Vault",
		Sources: cli.EnvVars("ANSIBLE_VAULT_ID", "INPUT_VAULT_ID", "PLUGIN_VAULT_ID"),
	},
	&cli.StringFlag{
		Name:    "vault-password",
		Usage:   "Password for decrypting an Ansible Vault",
		Sources: cli.EnvVars("ANSIBLE_VAULT_PASSWORD", "INPUT_VAULT_PASSWORD", "PLUGIN_VAULT_PASSWORD"),
	},
	&cli.IntFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		Usage:   "Set the verbosity level, ranging from 0 (minimal output) to 4 (maximum verbosity)",
		Sources: cli.EnvVars("ANSIBLE_VERBOSE", "INPUT_VERBOSE", "PLUGIN_VERBOSE"),
	},
	&cli.StringFlag{
		Name:    "private-key",
		Aliases: []string{"k"},
		Usage:   "Path to the SSH private key for connection",
		Sources: cli.EnvVars("ANSIBLE_PRIVATE_KEY", "INPUT_PRIVATE_KEY", "PLUGIN_PRIVATE_KEY"),
	},
	&cli.StringFlag{
		Name:    "private-key-file",
		Usage:   "Path to the file containing the SSH private key",
		Sources: cli.EnvVars("ANSIBLE_PRIVATE_KEY_FILE", "INPUT_PRIVATE_KEY_FILE", "PLUGIN_PRIVATE_KEY_FILE"),
	},
	&cli.StringFlag{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "Username to use for the connection",
		Sources: cli.EnvVars("ANSIBLE_USER", "INPUT_USER", "PLUGIN_USER"),
	},
	&cli.StringFlag{
		Name:    "connection",
		Aliases: []string{"c"},
		Usage:   "Type of connection to use (e.g., SSH)",
		Sources: cli.EnvVars("ANSIBLE_CONNECTION", "INPUT_CONNECTION", "PLUGIN_CONNECTION"),
	},
	&cli.IntFlag{
		Name:    "timeout",
		Aliases: []string{"T"},
		Usage:   "Override the connection timeout (in seconds)",
		Sources: cli.EnvVars("ANSIBLE_TIMEOUT", "INPUT_TIMEOUT", "PLUGIN_TIMEOUT"),
	},
	&cli.StringFlag{
		Name:    "ssh-common-args",
		Usage:   "Common arguments passed to all SSH-based connection methods",
		Sources: cli.EnvVars("ANSIBLE_SSH_COMMON_ARGS", "INPUT_SSH_COMMON_ARGS", "PLUGIN_SSH_COMMON_ARGS"),
	},
	&cli.StringFlag{
		Name:    "sftp-extra-args",
		Usage:   "Extra arguments passed exclusively to SFTP",
		Sources: cli.EnvVars("ANSIBLE_SFTP_EXTRA_ARGS", "INPUT_SFTP_EXTRA_ARGS", "PLUGIN_SFTP_EXTRA_ARGS"),
	},
	&cli.StringFlag{
		Name:    "scp-extra-args",
		Usage:   "Extra arguments passed exclusively to SCP",
		Sources: cli.EnvVars("ANSIBLE_SCP_EXTRA_ARGS", "INPUT_SCP_EXTRA_ARGS", "PLUGIN_SCP_EXTRA_ARGS"),
	},
	&cli.StringFlag{
		Name:    "ssh-extra-args",
		Usage:   "Extra arguments passed exclusively to SSH",
		Sources: cli.EnvVars("ANSIBLE_SSH_EXTRA_ARGS", "INPUT_SSH_EXTRA_ARGS", "PLUGIN_SSH_EXTRA_ARGS"),
	},
	&cli.BoolFlag{
		Name:    "become",
		Aliases: []string{"b"},
		Usage:   "Enable privilege escalation to run tasks as another user",
		Sources: cli.EnvVars("ANSIBLE_BECOME", "INPUT_BECOME", "PLUGIN_BECOME"),
	},
	&cli.StringFlag{
		Name:    "become-method",
		Usage:   "Method to use for privilege escalation (e.g., sudo)",
		Sources: cli.EnvVars("ANSIBLE_BECOME_METHOD", "INPUT_BECOME_METHOD", "PLUGIN_BECOME_METHOD"),
	},
	&cli.StringFlag{
		Name:    "become-user",
		Usage:   "User to impersonate when using privilege escalation",
		Sources: cli.EnvVars("ANSIBLE_BECOME_USER", "INPUT_BECOME_USER", "PLUGIN_BECOME_USER"),
	},
	&cli.BoolFlag{
		Name:    "ask-become-pass",
		Usage:   "Prompt for the become password",
		Sources: cli.EnvVars("ANSIBLE_ASK_BECOME_PASS", "INPUT_ASK_BECOME_PASS", "PLUGIN_ASK_BECOME_PASS"),
	},
	&cli.BoolFlag{
		Name:    "ask-pass",
		Usage:   "Prompt for the SSH password",
		Sources: cli.EnvVars("ANSIBLE_ASK_PASS", "INPUT_ASK_PASS", "PLUGIN_ASK_PASS"),
	},
	&cli.BoolFlag{
		Name:    "step",
		Usage:   "Prompt for confirmation before each task",
		Sources: cli.EnvVars("ANSIBLE_STEP", "INPUT_STEP", "PLUGIN_STEP"),
	},
	&cli.StringFlag{
		Name:    "ssh-transfer-method",
		Usage:   "Method for file transfer over SSH (e.g., scp or sftp)",
		Sources: cli.EnvVars("ANSIBLE_SSH_TRANSFER_METHOD", "INPUT_SSH_TRANSFER_METHOD", "PLUGIN_SSH_TRANSFER_METHOD"),
	},
	&cli.StringFlag{
		Name:    "output-callback",
		Usage:   "Set the stdout callback plugin for Ansible output",
		Sources: cli.EnvVars("ANSIBLE_OUTPUT_CALLBACK", "INPUT_OUTPUT_CALLBACK", "PLUGIN_OUTPUT_CALLBACK"),
	},
	&cli.BoolFlag{
		Name:    "no-color",
		Usage:   "Disable colorized output",
		Sources: cli.EnvVars("ANSIBLE_NO_COLOR", "INPUT_NO_COLOR", "PLUGIN_NO_COLOR"),
	},
	&cli.StringFlag{
		Name:    "vault-password-file",
		Usage:   "Path to a file containing the vault password",
		Sources: cli.EnvVars("ANSIBLE_VAULT_PASSWORD_FILE", "INPUT_VAULT_PASSWORD_FILE", "PLUGIN_VAULT_PASSWORD_FILE"),
	},
	&cli.BoolFlag{
		Name:    "ask-vault-pass",
		Usage:   "Prompt for the vault password",
		Sources: cli.EnvVars("ANSIBLE_ASK_VAULT_PASS", "INPUT_ASK_VAULT_PASS", "PLUGIN_ASK_VAULT_PASS"),
	},
	&cli.StringFlag{
		Name:    "fact-path",
		Usage:   "Path to local fact files",
		Sources: cli.EnvVars("ANSIBLE_FACT_PATH", "INPUT_FACT_PATH", "PLUGIN_FACT_PATH"),
	},
	&cli.StringFlag{
		Name:    "fact-caching",
		Usage:   "Caching method to use for facts",
		Sources: cli.EnvVars("ANSIBLE_FACT_CACHING", "INPUT_FACT_CACHING", "PLUGIN_FACT_CACHING"),
	},
	&cli.IntFlag{
		Name:    "fact-caching-timeout",
		Usage:   "Timeout (in seconds) for fact caching",
		Sources: cli.EnvVars("ANSIBLE_FACT_CACHING_TIMEOUT", "INPUT_FACT_CACHING_TIMEOUT", "PLUGIN_FACT_CACHING_TIMEOUT"),
	},
	&cli.StringFlag{
		Name:  "callbacks-enabled",
		Usage: "Comma-separated list of enabled callback plugins",
		Sources: cli.EnvVars(
			"ANSIBLE_CALLBACKS_ENABLED", "INPUT_CALLBACKS_ENABLED", "PLUGIN_CALLBACKS_ENABLED",
			// deprecated aliases - remove in next major version
			"ANSIBLE_CALLBACK_WHITELIST", "INPUT_CALLBACK_WHITELIST", "PLUGIN_CALLBACK_WHITELIST",
		),
	},
	&cli.IntFlag{
		Name:    "poll-interval",
		Usage:   "Interval (in seconds) for polling",
		Sources: cli.EnvVars("ANSIBLE_POLL_INTERVAL", "INPUT_POLL_INTERVAL", "PLUGIN_POLL_INTERVAL"),
	},
	&cli.StringFlag{
		Name:    "gather-subset",
		Usage:   "Limit the scope of gathered facts",
		Sources: cli.EnvVars("ANSIBLE_GATHER_SUBSET", "INPUT_GATHER_SUBSET", "PLUGIN_GATHER_SUBSET"),
	},
	&cli.IntFlag{
		Name:    "gather-timeout",
		Usage:   "Timeout (in seconds) for gathering facts",
		Sources: cli.EnvVars("ANSIBLE_GATHER_TIMEOUT", "INPUT_GATHER_TIMEOUT", "PLUGIN_GATHER_TIMEOUT"),
	},
	&cli.StringFlag{
		Name:    "strategy-plugin",
		Usage:   "Specify the strategy plugin to use",
		Sources: cli.EnvVars("ANSIBLE_STRATEGY_PLUGIN", "INPUT_STRATEGY_PLUGIN", "PLUGIN_STRATEGY_PLUGIN"),
	},
	&cli.IntFlag{
		Name:    "max-fail-percentage",
		Usage:   "Max percentage of hosts that can fail before the playbook aborts",
		Sources: cli.EnvVars("ANSIBLE_MAX_FAIL_PERCENTAGE", "INPUT_MAX_FAIL_PERCENTAGE", "PLUGIN_MAX_FAIL_PERCENTAGE"),
	},
	&cli.BoolFlag{
		Name:    "any-errors-fatal",
		Usage:   "Treat any error as fatal",
		Sources: cli.EnvVars("ANSIBLE_ANY_ERRORS_FATAL", "INPUT_ANY_ERRORS_FATAL", "PLUGIN_ANY_ERRORS_FATAL"),
	},
	&cli.StringFlag{
		Name:    "config-file",
		Usage:   "Path to the configuration file",
		Sources: cli.EnvVars("ANSIBLE_CONFIG_FILE", "INPUT_CONFIG_FILE", "PLUGIN_CONFIG_FILE"),
	},
	&cli.StringFlag{
		Name:    "temp-dir",
		Usage:   "Directory for temporary files",
		Sources: cli.EnvVars("ANSIBLE_TEMP_DIR", "INPUT_TEMP_DIR", "PLUGIN_TEMP_DIR"),
	},
}

func main() {
	// Load environment file if specified.
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		if err := godotenv.Load(filename); err != nil {
			log.Printf("Warning: Could not load env file: %v", err)
		}
	}

	cmd := &cli.Command{
		Name:  "Ansible Playbook Wrapper",
		Usage: "Execute Ansible Playbooks",
		Authors: []any{
			"arillso <hello@arillso.io>",
		},
		Flags:  appFlags,
		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		var ansibleErr *ansible.AnsibleError
		if errors.As(err, &ansibleErr) {
			if ansibleErr.ExitCode != 0 {
				log.Printf("Error: %v", err)
			}
			os.Exit(ansibleErr.ExitCode)
		}
		log.Fatalf("Error: %v", err)
	}
}

// validateParameters checks parameter integrity before execution.
func validateParameters(c *cli.Command) error {
	// Validate that required inventory files exist.
	for _, inv := range c.StringSlice("inventory") {
		if _, err := os.Stat(inv); os.IsNotExist(err) {
			return fmt.Errorf("%w: inventory file does not exist: %s", ErrInvalidParameter, inv)
		}
	}

	// Validate that required playbook files exist.
	for _, pb := range c.StringSlice("playbook") {
		if _, err := os.Stat(pb); os.IsNotExist(err) {
			return fmt.Errorf("%w: playbook file does not exist: %s", ErrInvalidParameter, pb)
		}
	}

	// Validate Galaxy file if specified.
	if galaxyFile := c.String("galaxy-file"); galaxyFile != "" {
		if _, err := os.Stat(galaxyFile); os.IsNotExist(err) {
			return fmt.Errorf("%w: galaxy file does not exist: %s", ErrInvalidParameter, galaxyFile)
		}
	}

	return nil
}

// sshAgent holds the state of a running ssh-agent process.
type sshAgent struct {
	sock string
	pid  string
}

// startSSHAgent starts an ssh-agent, adds the given private key, and returns
// the agent state for cleanup. The key content is written to a temporary file
// which is removed immediately after being added to the agent.
func startSSHAgent(privateKey, passphrase string) (*sshAgent, error) {
	// Start ssh-agent and capture its output.
	var buf bytes.Buffer
	agentCmd := exec.Command("ssh-agent", "-s")
	agentCmd.Stdout = &buf
	if err := agentCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to start ssh-agent: %w", err)
	}

	// Parse SSH_AUTH_SOCK and SSH_AGENT_PID from agent output.
	output := buf.String()
	agent := &sshAgent{}
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "SSH_AUTH_SOCK=") {
			agent.sock = strings.SplitN(strings.TrimPrefix(line, "SSH_AUTH_SOCK="), ";", 2)[0]
		}
		if strings.HasPrefix(line, "SSH_AGENT_PID=") {
			agent.pid = strings.SplitN(strings.TrimPrefix(line, "SSH_AGENT_PID="), ";", 2)[0]
		}
	}

	if agent.sock == "" {
		return nil, fmt.Errorf("failed to parse SSH_AUTH_SOCK from ssh-agent output")
	}

	// Write private key to a temporary file for ssh-add.
	tmpFile, err := os.CreateTemp("", "ssh-agent-key-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp key file: %w", err)
	}
	keyPath := tmpFile.Name()

	// Normalize line endings and ensure trailing newline.
	normalized := strings.ReplaceAll(privateKey, "\r\n", "\n")
	if !strings.HasSuffix(normalized, "\n") {
		normalized += "\n"
	}

	if _, err := tmpFile.WriteString(normalized); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(keyPath)
		return nil, fmt.Errorf("failed to write temp key file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(keyPath)
		return nil, fmt.Errorf("failed to close temp key file: %w", err)
	}

	// Add the key to the agent.
	addCtx, addCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer addCancel()
	addCmd := exec.CommandContext(addCtx, "ssh-add", keyPath)
	addCmd.Env = append(os.Environ(), "SSH_AUTH_SOCK="+agent.sock)

	// If a passphrase is provided, use SSH_ASKPASS to supply it non-interactively.
	if passphrase != "" {
		askpassScript, err := os.CreateTemp("", "ssh-askpass-")
		if err != nil {
			_ = os.Remove(keyPath)
			return nil, fmt.Errorf("failed to create askpass script: %w", err)
		}
		askpassPath := askpassScript.Name()
		// The script outputs the passphrase once, then exits with error on retries
		// to prevent ssh-add from looping indefinitely on wrong passphrases.
		flagFile := askpassPath + ".used"
		scriptContent := fmt.Sprintf("#!/bin/sh\nif [ -f '%s' ]; then exit 1; fi\ntouch '%s'\necho '%s'\n",
			flagFile, flagFile, strings.ReplaceAll(passphrase, "'", "'\\''"))
		if _, err := askpassScript.WriteString(scriptContent); err != nil {
			_ = askpassScript.Close()
			_ = os.Remove(askpassPath)
			_ = os.Remove(keyPath)
			return nil, fmt.Errorf("failed to write askpass script: %w", err)
		}
		_ = askpassScript.Close()
		if err := os.Chmod(askpassPath, 0700); err != nil {
			_ = os.Remove(askpassPath)
			_ = os.Remove(keyPath)
			return nil, fmt.Errorf("failed to set askpass script permissions: %w", err)
		}
		addCmd.Env = append(addCmd.Env,
			"SSH_ASKPASS="+askpassPath,
			"SSH_ASKPASS_REQUIRE=force",
			"DISPLAY=",
		)
		// Detach from controlling TTY so ssh-add uses SSH_ASKPASS.
		addCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		addCmd.Stdin = nil
		defer func() {
			_ = os.Remove(askpassPath)
			_ = os.Remove(flagFile)
		}()
	}

	if out, err := addCmd.CombinedOutput(); err != nil {
		_ = os.Remove(keyPath)
		return nil, fmt.Errorf("failed to add key to ssh-agent: %w: %s", err, string(out))
	}

	// Remove the temporary key file immediately after adding to agent.
	_ = os.Remove(keyPath)

	log.Printf("SSH agent started and key added (PID: %s)", agent.pid)
	return agent, nil
}

// stop kills the ssh-agent process.
func (a *sshAgent) stop() {
	if a == nil || a.pid == "" {
		return
	}
	killCmd := exec.Command("ssh-agent", "-k")
	killCmd.Env = append(os.Environ(), "SSH_AGENT_PID="+a.pid, "SSH_AUTH_SOCK="+a.sock)
	if err := killCmd.Run(); err != nil {
		log.Printf("Warning: failed to stop ssh-agent (PID %s): %v", a.pid, err)
	} else {
		log.Printf("SSH agent stopped (PID: %s)", a.pid)
	}
}

// run is the main action for executing the playbooks.
func run(ctx context.Context, c *cli.Command) error {
	// Validate parameters.
	if err := validateParameters(c); err != nil {
		return err
	}

	// Set execution timeout based on flag.
	timeoutDuration := time.Duration(c.Int("execution-timeout")) * time.Minute
	log.Printf("Setting execution timeout to %v minutes", c.Int("execution-timeout"))

	// Create context with timeout.
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	// Start ssh-agent if a private key is provided, so that ProxyCommand
	// and bastion host connections also have access to the key.
	extraEnv := make(map[string]string)
	if privateKey := c.String("private-key"); privateKey != "" {
		agent, err := startSSHAgent(privateKey, c.String("private-key-passphrase"))
		if err != nil {
			return fmt.Errorf("could not start ssh-agent for private key: %w", err)
		}
		defer agent.stop()
		extraEnv["SSH_AUTH_SOCK"] = agent.sock
	}

	log.Printf("Starting Ansible playbook execution with %d playbooks", len(c.StringSlice("playbook")))

	playbook := &ansible.Playbook{
		Config: ansible.Config{
			// Galaxy-related configuration.
			GalaxyFile:                        c.String("galaxy-file"),
			GalaxyForce:                       c.Bool("galaxy-force"),
			GalaxyAPIKey:                      c.String("galaxy-api-key"),
			GalaxyAPIServerURL:                c.String("galaxy-api-server-url"),
			GalaxyCollectionsPath:             c.String("galaxy-collections-path"),
			GalaxyDisableGPGVerify:            c.Bool("galaxy-disable-gpg-verify"),
			GalaxyForceWithDeps:               c.Bool("galaxy-force-with-deps"),
			GalaxyIgnoreCerts:                 c.Bool("galaxy-ignore-certs"),
			GalaxyIgnoreSignatureStatusCodes:  c.StringSlice("galaxy-ignore-signature-status-codes"),
			GalaxyKeyring:                     c.String("galaxy-keyring"),
			GalaxyOffline:                     c.Bool("galaxy-offline"),
			GalaxyPre:                         c.Bool("galaxy-pre"),
			GalaxyRequiredValidSignatureCount: c.Int("galaxy-required-valid-signature-count"),
			GalaxyRequirementsFile:            c.String("galaxy-requirements-file"),
			GalaxySignature:                   c.String("galaxy-signature"),
			GalaxyTimeout:                     c.Int("galaxy-timeout"),
			GalaxyUpgrade:                     c.Bool("galaxy-upgrade"),
			GalaxyNoDeps:                      c.Bool("galaxy-no-deps"),

			// Inventory and playbook configuration.
			Inventories:   c.StringSlice("inventory"),
			Playbooks:     c.StringSlice("playbook"),
			Limit:         c.String("limit"),
			SkipTags:      c.String("skip-tags"),
			StartAtTask:   c.String("start-at-task"),
			Tags:          c.String("tags"),
			ExtraVars:     c.StringSlice("extra-vars"),
			ModulePath:    c.StringSlice("module-path"),
			Check:         c.Bool("check"),
			Diff:          c.Bool("diff"),
			FlushCache:    c.Bool("flush-cache"),
			ForceHandlers: c.Bool("force-handlers"),
			ListHosts:     c.Bool("list-hosts"),
			ListTags:      c.Bool("list-tags"),
			ListTasks:     c.Bool("list-tasks"),
			SyntaxCheck:   c.Bool("syntax-check"),
			Forks:         c.Int("forks"),

			// Vault and authentication configuration.
			VaultID:            c.String("vault-id"),
			VaultPassword:      c.String("vault-password"),
			Verbose:            c.Int("verbose"),
			PrivateKey:         c.String("private-key"),
			PrivateKeyFile:     c.String("private-key-file"),
			User:               c.String("user"),
			Connection:         c.String("connection"),
			Timeout:            c.Int("timeout"),
			SSHCommonArgs:      c.String("ssh-common-args"),
			SCPExtraArgs:       c.String("scp-extra-args"),
			SFTPExtraArgs:      c.String("sftp-extra-args"),
			SSHExtraArgs:       c.String("ssh-extra-args"),
			Become:             c.Bool("become"),
			BecomeMethod:       c.String("become-method"),
			BecomeUser:         c.String("become-user"),
			AskBecomePass:      c.Bool("ask-become-pass"),
			AskPass:            c.Bool("ask-pass"),
			Step:               c.Bool("step"),
			SSHTransferMethod:  c.String("ssh-transfer-method"),
			NoColor:            c.Bool("no-color"),
			OutputCallback:     c.String("output-callback"),
			VaultPasswordFile:  c.String("vault-password-file"),
			AskVaultPass:       c.Bool("ask-vault-pass"),
			FactPath:           c.String("fact-path"),
			FactCaching:        c.String("fact-caching"),
			FactCachingTimeout: c.Int("fact-caching-timeout"),
			CallbacksEnabled:   c.String("callbacks-enabled"),
			PollInterval:       c.Int("poll-interval"),
			GatherSubset:       c.String("gather-subset"),
			GatherTimeout:      c.Int("gather-timeout"),
			StrategyPlugin:     c.String("strategy-plugin"),
			MaxFailPercentage:  c.Int("max-fail-percentage"),
			AnyErrorsFatal:     c.Bool("any-errors-fatal"),
			ConfigFile:         c.String("config-file"),
			TempDir:            c.String("temp-dir"),
			ExtraEnv:           extraEnv,
		},
	}

	return playbook.Exec(ctx)
}
