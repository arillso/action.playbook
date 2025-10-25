# Action: Play Ansible Playbook

Github Action for running Ansible Playbooks.

## Inputs

### galaxy_file

Specifies the path to the Ansible Galaxy requirements file.

### galaxy_force

Forces the reinstallation of roles or collections from the Galaxy file.

### galaxy_api_key

Sets the API key used for authenticating to Ansible Galaxy.

### galaxy_api_server_url

Defines the URL of the Ansible Galaxy API server to interact with.

### galaxy_collections_path

Sets the path to the directory where Galaxy collections are stored.

### galaxy_disable_gpg_verify

Disables GPG signature verification for Ansible Galaxy operations.

### galaxy_force_with_deps

Forces the installation of collections with their dependencies from Galaxy.

### galaxy_ignore_certs

Ignores SSL certificate validation for Ansible Galaxy requests.

### galaxy_ignore_signature_status_codes

Lists HTTP status codes to ignore during Galaxy signature validation.

### galaxy_keyring

Specifies the path to the GPG keyring used with Ansible Galaxy.

### galaxy_offline

Enables offline mode, preventing any requests to Ansible Galaxy.

### galaxy_pre

Allows the installation of pre-release versions from Ansible Galaxy.

### galaxy_required_valid_signature_count

Sets the required number of valid GPG signatures for Galaxy content.

### galaxy_requirements_file

Defines the path to the Ansible Galaxy requirements file.

### galaxy_signature

Specifies a specific GPG signature to verify for Galaxy content.

### galaxy_timeout

Sets the timeout in seconds for Ansible Galaxy operations.

### galaxy_upgrade

Enables automatic upgrading of Galaxy collections to the latest version.

### galaxy_no_deps

Disables automatic resolution of dependencies in Ansible Galaxy.

### inventory

**Required.** Specifies one or more inventory host files for Ansible to use.

### playbook

**Required.** List of playbooks to apply.

### limit

Further limit selected hosts to an additional pattern.

### skip_tags

Only run plays and tasks whose tags do not match these values.

### start_at_task

Start the playbook at the task matching this name.

### tags

Only run plays and tasks tagged with these values.

### extra_vars

Set additional variables in a key=value format for the playbook.

### module_path

Prepends specified paths to the module library path list.

### check

Executes a dry run, showing what changes would be made without making them.

### diff

Shows the differences in files and templates when changing them.

### flush_cache

Clears the fact cache for every host in the inventory.

### force_handlers

Runs all handlers even if a task fails.

### list_hosts

Outputs a list of matching hosts.

### list_tags

List all available tags.

### list_tasks

List all tasks that would be executed.

### syntax_check

Performs a syntax check on the playbook, without executing it.

### forks

Defines the number of parallel processes to use during playbook execution.

### vault_id

Specifies the identity to use when accessing an Ansible Vault.

### vault_password

The vault password to use. This should be stored in a Secret on Github.

### verbose

Sets the verbosity level, ranging from 0 (minimal output) to 4 (maximum verbosity).

### private_key

Specifies the SSH private key content for connections. The key will be written to a temporary file by the action. This should be stored in a Secret on Github.

**Starting with v1.2.0:** The action automatically normalizes SSH keys (converts CRLF to LF, adds trailing newlines, validates format).

**Note:** If you're experiencing issues with bastion hosts or ProxyCommand configurations on older versions, see the [SSH Authentication with Bastion Hosts](#ssh-authentication-with-bastion-hosts) section for workarounds.

### user

Defines the username for making connections.

### connection

Sets the type of connection to use (e.g., SSH).

### timeout

Overrides the default connection timeout in seconds.

### ssh_common_args

Specifies common arguments to pass to all SSH-based connection methods (SSH, SCP, SFTP).

### sftp_extra_args

Provides extra arguments to pass only to SFTP.

### scp_extra_args

Provides extra arguments to pass only to SCP.

### ssh_extra_args

Provides extra arguments to pass only to SSH.

### become

Enables privilege escalation, allowing operations to run as another user.

### become_method

Specifies the method to use for privilege escalation (e.g., sudo).

### become_user

Sets the user to impersonate when using privilege escalation.

## SSH Authentication

> **✨ New in v1.2.0+**
>
> SSH private keys are automatically normalized:
> - Windows line endings (CRLF) → Unix format (LF)
> - Missing trailing newlines are added
> - Format validation (RSA, OpenSSH, EC, DSA)

### Basic Usage

```yaml
- name: Deploy with Ansible
  uses: arillso/action.playbook@1.2.0
  with:
    playbook: deploy.yml
    inventory: ansible_hosts.yml
    private_key: ${{ secrets.SSH_PRIVATE_KEY }}
  env:
    ANSIBLE_HOST_KEY_CHECKING: 'false'
```

### Storing SSH Keys in GitHub Secrets

> **🔒 Security Warning: NEVER commit private SSH keys to your repository!**

1. Go to repository **Settings** → **Secrets and variables** → **Actions**
2. Click **"New repository secret"**
3. Name: `SSH_PRIVATE_KEY`
4. Value: Paste your entire private key (including `-----BEGIN/END-----` headers)
5. Click **"Add secret"**

Your key must have proper PEM headers (`-----BEGIN RSA PRIVATE KEY-----` or similar). Line endings and trailing newlines are automatically fixed in v1.2.0+.

### Troubleshooting

**For versions < 1.2.0** or edge cases, use `extra_vars` to pass the key file path:

```yaml
- name: Setup SSH Key
  run: |
    echo "${{ secrets.SSH_PRIVATE_KEY }}" > /tmp/ssh_key
    chmod 600 /tmp/ssh_key

- name: Deploy
  uses: arillso/action.playbook@master
  with:
    playbook: deploy.yml
    inventory: ansible_hosts.yml
    extra_vars: ansible_ssh_private_key_file=/tmp/ssh_key
```

For verbose SSH debugging, set `verbose: 4` in the action inputs.

## Example Usage

```yaml
- name: Play Ansible Playbook
  uses: arillso/action.playbook@master
  with:
    playbook: tests/playbook.yml
    inventory: tests/hosts.yml
    galaxy_file: tests/requirements.yml
  env:
    ANSIBLE_HOST_KEY_CHECKING: 'false'
    ANSIBLE_DEPRECATION_WARNINGS: 'false'
```

<!-- ALL-CONTRIBUTORS-LIST:END -->

## License

<!-- markdownlint-disable -->

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.

<!-- markdownlint-enable -->

## Copyright

(c) 2020, Arillso
