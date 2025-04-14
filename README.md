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

Use this key to authenticate the connection. This should be stored in a Secret on Github.

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
