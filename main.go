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
		Copyright: "Copyright (c) 2020 Arillso",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "arillso",
				Email: "hello@arillso.io",
			},
		},
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "galaxy-file",
				Usage:   "path to galaxy requirements",
				EnvVars: []string{"ANSIBLE_GALAXY_FILE", "INPUT_GALAXY_FILE", "PLUGIN_GALAXY_FILE"},
			},
			&cli.StringSliceFlag{
				Name:     "inventory,i",
				Usage:    "specify inventory host path",
				EnvVars:  []string{"ANSIBLE_INVENTORY", "INPUT_INVENTORY", "PLUGIN_INVENTORY"},
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "playbook",
				Usage:    "list of playbooks to apply",
				EnvVars:  []string{"ANSIBLE_PLAYBOOK", "INPUT_PLAYBOOK", "PLUGIN_PLAYBOOK"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "limit,l",
				Usage:   "further limit selected hosts to an additional pattern",
				EnvVars: []string{"ANSIBLE_LIMIT", "INPUT_LIMIT", "PLUGIN_LIMIT"},
			},
			&cli.StringFlag{
				Name:    "skip-tags",
				Usage:   "only run plays and tasks whose tags do not match these values",
				EnvVars: []string{"ANSIBLE_SKIP_TAGS", "INPUT_SKIP_TAGS", "PLUGIN_SKIP_TAGS"},
			},
			&cli.StringFlag{
				Name:    "start-at-task",
				Usage:   "start the playbook at the task matching this name",
				EnvVars: []string{"ANSIBLE_START_AT_TASK", "INPUT_START_AT_TASK", "PLUGIN_START_AT_TASK"},
			},
			&cli.StringFlag{
				Name:    "tags,t",
				Usage:   "only run plays and tasks tagged with these values",
				EnvVars: []string{"ANSIBLE_TAGS", "INPUT_TAGS", "PLUGIN_TAGS"},
			},
			&cli.StringSliceFlag{
				Name:    "extra-vars",
				Usage:   "set additional variables as key=value",
				EnvVars: []string{"ANSIBLE_EXTRA_VARS", "INPUT_EXTRA_VARS", "PLUGIN_EXTRA_VARS"},
			},
			&cli.StringSliceFlag{
				Name:    "module-path,M",
				Usage:   "prepend paths to module library",
				EnvVars: []string{"ANSIBLE_MODULE_PATH", "INPUT_MODULE_PATH", "PLUGIN_MODULE_PATH"},
			},
			&cli.BoolFlag{
				Name:    "check,C",
				Usage:   "run a check, do not apply any changes",
				EnvVars: []string{"ANSIBLE_CHECK", "INPUT_CHECK", "PLUGIN_CHECK"},
			},
			&cli.BoolFlag{
				Name:    "diff,D",
				Usage:   "when changing (small) files and templates, show the differences in those files; works great with –check",
				EnvVars: []string{"ANSIBLE_DIFF", "INPUT_DIFF", "PLUGIN_DIFF"},
			},
			&cli.BoolFlag{
				Name:    "flush-cache",
				Usage:   "clear the fact cache for every host in inventory",
				EnvVars: []string{"ANSIBLE_FLUSH_CACHE", "INPUT_FLUSH_CACHE", "PLUGIN_FLUSH_CACHE"},
			},
			&cli.BoolFlag{
				Name:    "force-handlers",
				Usage:   "run handlers even if a task fails",
				EnvVars: []string{"ANSIBLE_FORCE_HANDLERS", "INPUT_FORCE_HANDLERS", "PLUGIN_FORCE_HANDLERS"},
			},
			&cli.BoolFlag{
				Name:    "list-hosts",
				Usage:   "outputs a list of matching hosts",
				EnvVars: []string{"ANSIBLE_LIST_HOSTS", "INPUT_LIST_HOSTS", "PLUGIN_LIST_HOSTS"},
			},
			&cli.BoolFlag{
				Name:    "list-tags",
				Usage:   "list all available tags",
				EnvVars: []string{"ANSIBLE_LIST_TAGS", "INPUT_LIST_TAGS", "PLUGIN_LIST_TAGS"},
			},
			&cli.BoolFlag{
				Name:    "list-tasks",
				Usage:   "list all tasks that would be executed",
				EnvVars: []string{"ANSIBLE_LIST_TASKS", "INPUT_LIST_TASKS", "PLUGIN_LIST_TASKS"},
			},
			&cli.BoolFlag{
				Name:    "syntax-check",
				Usage:   "perform a syntax check on the playbook",
				EnvVars: []string{"ANSIBLE_SYNTAX_CHECK", "INPUT_SYNTAX_CHECK", "PLUGIN_SYNTAX_CHECK"},
			},
			&cli.IntFlag{
				Name:    "forks,f",
				Usage:   "specify number of parallel processes to use",
				EnvVars: []string{"ANSIBLE_FORKS", "INPUT_FORKS", "PLUGIN_FORKS"},
				Value:   5,
			},
			&cli.StringFlag{
				Name:    "vault-id",
				Usage:   "the vault identity to use",
				EnvVars: []string{"ANSIBLE_VAULT_ID", "INPUT_VAULT_ID", "PLUGIN_VAULT_ID"},
			},
			&cli.StringFlag{
				Name:    "vault-password",
				Usage:   "the vault password to use",
				EnvVars: []string{"ANSIBLE_VAULT_PASSWORD", "INPUT_VAULT_PASSWORD", "PLUGIN_VAULT_PASSWORD"},
			},
			&cli.IntFlag{
				Name:    "verbose",
				Usage:   "level of verbosity, 0 up to 4",
				EnvVars: []string{"ANSIBLE_VERBOSE", "INPUT_VERBOSE", "PLUGIN_VERBOSE"},
			},
			&cli.StringFlag{
				Name:    "private-key",
				Usage:   "use this key to authenticate the connection",
				EnvVars: []string{"ANSIBLE_PRIVATE_KEY", "INPUT_PRIVATE_KEY", "PLUGIN_PRIVATE_KEY"},
			},
			&cli.StringFlag{
				Name:    "user",
				Usage:   "connect as this user",
				EnvVars: []string{"ANSIBLE_USER", "INPUT_USER", "PLUGIN_USER"},
			},
			&cli.StringFlag{
				Name:    "connection",
				Usage:   "connection type to use",
				EnvVars: []string{"ANSIBLE_CONNECTION", "INPUT_CONNECTION", "PLUGIN_CONNECTION"},
			},
			&cli.IntFlag{
				Name:    "timeout",
				Usage:   "override the connection timeout in seconds",
				EnvVars: []string{"ANSIBLE_TIMEOUT", "INPUT_TIMEOUT", "PLUGIN_TIMEOUT"},
			},
			&cli.StringFlag{
				Name:    "ssh-common-args",
				Usage:   "specify common arguments to pass to sftp/scp/ssh",
				EnvVars: []string{"ANSIBLE_SSH_COMMON_ARGS", "INPUT_SSH_COMMON_ARGS", "PLUGIN_SSH_COMMON_ARGS"},
			},
			&cli.StringFlag{
				Name:    "sftp-extra-args",
				Usage:   "specify extra arguments to pass to sftp only",
				EnvVars: []string{"ANSIBLE_SFTP_EXTRA_ARGS", "INPUT_SFTP_EXTRA_ARGS", "PLUGIN_SFTP_EXTRA_ARGS"},
			},
			&cli.StringFlag{
				Name:    "scp-extra-args",
				Usage:   "specify extra arguments to pass to scp only",
				EnvVars: []string{"ANSIBLE_SCP_EXTRA_ARGS", "INPUT_SCP_EXTRA_ARGS", "PLUGIN_SCP_EXTRA_ARGS"},
			},
			&cli.StringFlag{
				Name:    "ssh-extra-args",
				Usage:   "specify extra arguments to pass to ssh only",
				EnvVars: []string{"ANSIBLE_SSH_EXTRA_ARGS", "INPUT_SSH_EXTRA_ARGS", "PLUGIN_SSH_EXTRA_ARGS"},
			},
			&cli.BoolFlag{
				Name:    "become",
				Usage:   "run operations with become",
				EnvVars: []string{"ANSIBLE_BECOME", "INPUT_BECOME", "PLUGIN_BECOME"},
			},
			&cli.StringFlag{
				Name:    "become-method",
				Usage:   "privilege escalation method to use",
				EnvVars: []string{"ANSIBLE_BECOME_METHOD", "INPUT_BECOME_METHOD", "PLUGIN_BECOME_METHOD"},
			},
			&cli.StringFlag{
				Name:    "become-user",
				Usage:   "run operations as this user",
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
			GalaxyFile:    c.String("galaxy-file"),
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
			VaultID:       c.String("vailt-id"),
			VaultPassword: c.String("vault-password"),
			Verbose:       c.Int("verbose"),
			PrivateKey:    c.String("private-key"),
			User:          c.String("user"),
			Connection:    c.String("connection"),
			Timeout:       c.Int("timeout"),
			SSHCommonArgs: c.String("ssh-common-args"),
			SFTPExtraArgs: c.String("sftp-extra-args"),
			SCPExtraArgs:  c.String("scp-extra-args"),
			SSHExtraArgs:  c.String("ssh-extra-args"),
			Become:        c.Bool("become"),
			BecomeMethod:  c.String("become-method"),
			BecomeUser:    c.String("become-user"),
		},
	}

	return playbook.Exec()
}
