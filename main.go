package main

import (
	"log"
	"os"

	ansible "github.com/arillso/go.ansible"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version = "unknown"
)

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	app := &cli.App{
		Name:      "Ansible Playbook Wrapper",
		Usage:     "Executing Ansible Playbook",
		Copyright: "Copyright (c) 2023 Arillso",
		Authors: []*cli.Author{
			{
				Name:  "arillso",
				Email: "hello@arillso.io",
			},
		},
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "galaxy-file",
				Usage:   "Specifies the path to the Ansible Galaxy requirements file.",
				EnvVars: []string{"ANSIBLE_GALAXY_FILE", "INPUT_GALAXY_FILE", "PLUGIN_GALAXY_FILE"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-force",
				Usage:   "Forces the reinstallation of roles or collections from the Galaxy file.",
				EnvVars: []string{"ANSIBLE_GALAXY_FORCE", "INPUT_GALAXY_FORCE", "PLUGIN_GALAXY_FORCE"},
			},
			&cli.StringFlag{
				Name:    "galaxy-api-key",
				Usage:   "Sets the API key used for authenticating to Ansible Galaxy.",
				EnvVars: []string{"ANSIBLE_GALAXY_API_KEY"},
			},
			&cli.StringFlag{
				Name:    "galaxy-api-server-url",
				Usage:   "Defines the URL of the Ansible Galaxy API server to interact with.",
				EnvVars: []string{"ANSIBLE_GALAXY_API_SERVER_URL"},
			},
			&cli.StringFlag{
				Name:    "galaxy-collections-path",
				Usage:   "Sets the path to the directory where Galaxy collections are stored.",
				EnvVars: []string{"ANSIBLE_GALAXY_COLLECTIONS_PATH"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-disable-gpg-verify",
				Usage:   "Disables GPG signature verification for Ansible Galaxy operations.",
				EnvVars: []string{"ANSIBLE_GALAXY_DISABLE_GPG_VERIFY"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-force-with-deps",
				Usage:   "Forces the installation of collections with their dependencies from Galaxy.",
				EnvVars: []string{"ANSIBLE_GALAXY_FORCE_WITH_DEPS"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-ignore-certs",
				Usage:   "Ignores SSL certificate validation for Ansible Galaxy requests.",
				EnvVars: []string{"ANSIBLE_GALAXY_IGNORE_CERTS"},
			},
			&cli.StringSliceFlag{
				Name:    "galaxy-ignore-signature-status-codes",
				Usage:   "Lists HTTP status codes to ignore during Galaxy signature validation.",
				EnvVars: []string{"ANSIBLE_GALAXY_IGNORE_SIGNATURE_STATUS_CODES"},
			},
			&cli.StringFlag{
				Name:    "galaxy-keyring",
				Usage:   "Specifies the path to the GPG keyring used with Ansible Galaxy.",
				EnvVars: []string{"ANSIBLE_GALAXY_KEYRING"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-offline",
				Usage:   "Enables offline mode, preventing any requests to Ansible Galaxy.",
				EnvVars: []string{"ANSIBLE_GALAXY_OFFLINE"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-pre",
				Usage:   "Allows the installation of pre-release versions from Ansible Galaxy.",
				EnvVars: []string{"ANSIBLE_GALAXY_PRE"},
			},
			&cli.IntFlag{
				Name:    "galaxy-required-valid-signature-count",
				Usage:   "Sets the required number of valid GPG signatures for Galaxy content.",
				EnvVars: []string{"ANSIBLE_GALAXY_REQUIRED_VALID_SIGNATURE_COUNT"},
			},
			&cli.StringFlag{
				Name:    "galaxy-requirements-file",
				Usage:   "Defines the path to the Ansible Galaxy requirements file.",
				EnvVars: []string{"ANSIBLE_GALAXY_REQUIREMENTS_FILE"},
			},
			&cli.StringFlag{
				Name:    "galaxy-signature",
				Usage:   "Specifies a specific GPG signature to verify for Galaxy content.",
				EnvVars: []string{"ANSIBLE_GALAXY_SIGNATURE"},
			},
			&cli.IntFlag{
				Name:    "galaxy-timeout",
				Usage:   "Sets the timeout in seconds for Ansible Galaxy operations.",
				EnvVars: []string{"ANSIBLE_GALAXY_TIMEOUT"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-upgrade",
				Usage:   "Enables automatic upgrading of Galaxy collections to the latest version.",
				EnvVars: []string{"ANSIBLE_GALAXY_UPGRADE"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-no-deps",
				Usage:   "Disables automatic resolution of dependencies in Ansible Galaxy.",
				EnvVars: []string{"ANSIBLE_GALAXY_NO_DEPS"},
			},
			&cli.StringSliceFlag{
				Name:     "inventory",
				Aliases:  []string{"i"},
				Usage:    "Specifies one or more inventory host files for Ansible to use.",
				EnvVars:  []string{"ANSIBLE_INVENTORY", "INPUT_INVENTORY", "PLUGIN_INVENTORY"},
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "playbook",
				Aliases:  []string{"p"},
				Usage:    "List of playbooks to apply.",
				EnvVars:  []string{"ANSIBLE_PLAYBOOK", "INPUT_PLAYBOOK", "PLUGIN_PLAYBOOK"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Limits the playbook execution to a specific group of hosts.",
				EnvVars: []string{"ANSIBLE_LIMIT", "INPUT_LIMIT", "PLUGIN_LIMIT"},
			},
			&cli.StringFlag{
				Name:    "skip-tags",
				Usage:   "Only run plays and tasks whose tags do not match these values.",
				EnvVars: []string{"ANSIBLE_SKIP_TAGS", "INPUT_SKIP_TAGS", "PLUGIN_SKIP_TAGS"},
			},
			&cli.StringFlag{
				Name:    "start-at-task",
				Usage:   "Start the playbook at the task matching this name.",
				EnvVars: []string{"ANSIBLE_START_AT_TASK", "INPUT_START_AT_TASK", "PLUGIN_START_AT_TASK"},
			},
			&cli.StringFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "Executes only tasks and plays with specified tags.",
				EnvVars: []string{"ANSIBLE_TAGS", "INPUT_TAGS", "PLUGIN_TAGS"},
			},
			&cli.StringSliceFlag{
				Name:    "extra-vars",
				Aliases: []string{"e"},
				Usage:   "Sets additional variables in a key=value format for the playbook.",
				EnvVars: []string{"ANSIBLE_EXTRA_VARS", "INPUT_EXTRA_VARS", "PLUGIN_EXTRA_VARS"},
			},
			&cli.StringSliceFlag{
				Name:    "module-path",
				Aliases: []string{"M"},
				Usage:   "Prepends specified paths to the module library path list.",
				EnvVars: []string{"ANSIBLE_MODULE_PATH", "INPUT_MODULE_PATH", "PLUGIN_MODULE_PATH"},
			},
			&cli.BoolFlag{
				Name:    "check",
				Aliases: []string{"C"},
				Usage:   "Executes a dry run, showing what changes would be made without making them.",
				EnvVars: []string{"ANSIBLE_CHECK", "INPUT_CHECK", "PLUGIN_CHECK"},
			},
			&cli.BoolFlag{
				Name:    "diff",
				Aliases: []string{"D"},
				Usage:   "Shows the differences in files and templates when changing them.",
				EnvVars: []string{"ANSIBLE_DIFF", "INPUT_DIFF", "PLUGIN_DIFF"},
			},
			&cli.BoolFlag{
				Name:    "flush-cache",
				Usage:   "Clears the fact cache for every host in the inventory.",
				EnvVars: []string{"ANSIBLE_FLUSH_CACHE", "INPUT_FLUSH_CACHE", "PLUGIN_FLUSH_CACHE"},
			},
			&cli.BoolFlag{
				Name:    "force-handlers",
				Usage:   "Runs all handlers even if a task fails.",
				EnvVars: []string{"ANSIBLE_FORCE_HANDLERS", "INPUT_FORCE_HANDLERS", "PLUGIN_FORCE_HANDLERS"},
			},
			&cli.BoolFlag{
				Name:    "list-hosts",
				Usage:   "Outputs a list of matching hosts.",
				EnvVars: []string{"ANSIBLE_LIST_HOSTS", "INPUT_LIST_HOSTS", "PLUGIN_LIST_HOSTS"},
			},
			&cli.BoolFlag{
				Name:    "list-tags",
				Usage:   "List all available tags.",
				EnvVars: []string{"ANSIBLE_LIST_TAGS", "INPUT_LIST_TAGS", "PLUGIN_LIST_TAGS"},
			},
			&cli.BoolFlag{
				Name:    "list-tasks",
				Usage:   "List all tasks that would be executed.",
				EnvVars: []string{"ANSIBLE_LIST_TASKS", "INPUT_LIST_TASKS", "PLUGIN_LIST_TASKS"},
			},
			&cli.BoolFlag{
				Name:    "syntax-check",
				Usage:   "Performs a syntax check on the playbook, without executing it.",
				EnvVars: []string{"ANSIBLE_SYNTAX_CHECK", "INPUT_SYNTAX_CHECK", "PLUGIN_SYNTAX_CHECK"},
			},
			&cli.IntFlag{
				Name:    "forks",
				Aliases: []string{"f"},
				Usage:   "Defines the number of parallel processes to use during playbook execution.",
				EnvVars: []string{"ANSIBLE_FORKS", "INPUT_FORKS", "PLUGIN_FORKS"},
				Value:   5,
			},
			&cli.StringFlag{
				Name:    "vault-id",
				Usage:   "Specifies the identity to use when accessing an Ansible Vault.",
				EnvVars: []string{"ANSIBLE_VAULT_ID", "INPUT_VAULT_ID", "PLUGIN_VAULT_ID"},
			},

			&cli.StringFlag{
				Name:    "vault-password",
				Usage:   "Sets the password to use for decrypting an Ansible Vault.",
				EnvVars: []string{"ANSIBLE_VAULT_PASSWORD", "INPUT_VAULT_PASSWORD", "PLUGIN_VAULT_PASSWORD"},
			},
			&cli.IntFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Sets the verbosity level, ranging from 0 (minimal output) to 4 (maximum verbosity).",
				EnvVars: []string{"ANSIBLE_VERBOSE", "INPUT_VERBOSE", "PLUGIN_VERBOSE"},
			},
			&cli.StringFlag{
				Name:    "private-key",
				Aliases: []string{"k"},
				Usage:   "Specifies the SSH private key file for connections.",
				EnvVars: []string{"ANSIBLE_PRIVATE_KEY", "INPUT_PRIVATE_KEY", "PLUGIN_PRIVATE_KEY"},
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Defines the username for making connections.",
				EnvVars: []string{"ANSIBLE_USER", "INPUT_USER", "PLUGIN_USER"},
			},
			&cli.StringFlag{
				Name:    "connection",
				Aliases: []string{"c"},
				Usage:   "Sets the type of connection to use (e.g., SSH).",
				EnvVars: []string{"ANSIBLE_CONNECTION", "INPUT_CONNECTION", "PLUGIN_CONNECTION"},
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"T"},
				Usage:   "Overrides the default connection timeout in seconds.",
				EnvVars: []string{"ANSIBLE_TIMEOUT", "INPUT_TIMEOUT", "PLUGIN_TIMEOUT"},
			},
			&cli.StringFlag{
				Name:    "ssh-common-args",
				Usage:   "Specifies common arguments to pass to all SSH-based connection methods (SSH, SCP, SFTP).",
				EnvVars: []string{"ANSIBLE_SSH_COMMON_ARGS", "INPUT_SSH_COMMON_ARGS", "PLUGIN_SSH_COMMON_ARGS"},
			},
			&cli.StringFlag{
				Name:    "sftp-extra-args",
				Usage:   "Provides extra arguments to pass only to SFTP.",
				EnvVars: []string{"ANSIBLE_SFTP_EXTRA_ARGS", "INPUT_SFTP_EXTRA_ARGS", "PLUGIN_SFTP_EXTRA_ARGS"},
			},
			&cli.StringFlag{
				Name:    "scp-extra-args",
				Usage:   "Provides extra arguments to pass only to SCP.",
				EnvVars: []string{"ANSIBLE_SCP_EXTRA_ARGS", "INPUT_SCP_EXTRA_ARGS", "PLUGIN_SCP_EXTRA_ARGS"},
			},
			&cli.StringFlag{
				Name:    "ssh-extra-args",
				Usage:   "Provides extra arguments to pass only to SSH.",
				EnvVars: []string{"ANSIBLE_SSH_EXTRA_ARGS", "INPUT_SSH_EXTRA_ARGS", "PLUGIN_SSH_EXTRA_ARGS"},
			},
			&cli.BoolFlag{
				Name:    "become",
				Aliases: []string{"b"},
				Usage:   "Enables privilege escalation, allowing operations to run as another user.",
				EnvVars: []string{"ANSIBLE_BECOME", "INPUT_BECOME", "PLUGIN_BECOME"},
			},
			&cli.StringFlag{
				Name:    "become-method",
				Usage:   "Specifies the method to use for privilege escalation (e.g., sudo).",
				EnvVars: []string{"ANSIBLE_BECOME_METHOD", "INPUT_BECOME_METHOD", "PLUGIN_BECOME_METHOD"},
			},
			&cli.StringFlag{
				Name:    "become-user",
				Usage:   "Sets the user to impersonate when using privilege escalation.",
				EnvVars: []string{"ANSIBLE_BECOME_USER", "INPUT_BECOME_USER", "PLUGIN_BECOME_USER"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	playbook := &ansible.AnsiblePlaybook{
		Config: ansible.Config{
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
			Inventories:                       c.StringSlice("inventory"),
			Playbooks:                         c.StringSlice("playbook"),
			Limit:                             c.String("limit"),
			SkipTags:                          c.String("skip-tags"),
			StartAtTask:                       c.String("start-at-task"),
			Tags:                              c.String("tags"),
			ExtraVars:                         c.StringSlice("extra-vars"),
			ModulePath:                        c.StringSlice("module-path"),
			Check:                             c.Bool("check"),
			Diff:                              c.Bool("diff"),
			FlushCache:                        c.Bool("flush-cache"),
			ForceHandlers:                     c.Bool("force-handlers"),
			ListHosts:                         c.Bool("list-hosts"),
			ListTags:                          c.Bool("list-tags"),
			ListTasks:                         c.Bool("list-tasks"),
			SyntaxCheck:                       c.Bool("syntax-check"),
			Forks:                             c.Int("forks"),
			VaultID:                           c.String("vailt-id"),
			VaultPassword:                     c.String("vault-password"),
			Verbose:                           c.Int("verbose"),
			PrivateKey:                        c.String("private-key"),
			User:                              c.String("user"),
			Connection:                        c.String("connection"),
			Timeout:                           c.Int("timeout"),
			SSHCommonArgs:                     c.String("ssh-common-args"),
			SFTPExtraArgs:                     c.String("sftp-extra-args"),
			SCPExtraArgs:                      c.String("scp-extra-args"),
			SSHExtraArgs:                      c.String("ssh-extra-args"),
			Become:                            c.Bool("become"),
			BecomeMethod:                      c.String("become-method"),
			BecomeUser:                        c.String("become-user"),
		},
	}

	return playbook.Exec()
}
