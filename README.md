# Action: Ansible Playbook

Github Action for running Ansible Playbooks.

## Inputs

### galaxy_file

Name of the galaxy file in your workspace.

### inventory

**Required** inventory

Name of the inventory file in your workspace.

### playbook

**Required** playbook

Name of the playbook in your workspace.

### limit

Further limit selected hosts to an additional pattern.

### skip_tags

Only run plays and tasks whose tags do not match these values.

### start_at_task

Start the playbook at the task matching this name.

### tags

Only run plays and tasks tagged with these values.

### extra_vars

Set additional variables as key=value.

### module_path

Prepend paths to module library.

### check

Run a check, do not apply any changes.

### diff

When changing (small) files and templates, show the differences in those files; works great with –check.

### flush_cache

Clear the fact cache for every host in inventory.

### force_handlers

Run handlers even if a task fails.

### list_hosts

Outputs a list of matching hosts.

### list_tags

List all available tags.

### list_tasks

List all tasks that would be executed.

### syntax_check

Perform a syntax check on the playbook.

### forks

Specify number of parallel processes to use.

### vault_id

The vault identity to use.

### vault_password

The vault password to use. This should be stored in a Secret on Github.

See [https://help.github.com/en/github/automating-your-workflow-with-github-actions/virtual-environments-for-github-actions#creating-and-using-secrets-encrypted-variables](https://help.github.com/en/github/automating-your-workflow-with-github-actions/virtual-environments-for-github-actions#creating-and-using-secrets-encrypted-variables)

### verbose

Level of verbosity, 0 up to 4.

### private_key

Use this key to authenticate the connection. This should be stored in a Secret on Github.

See [https://help.github.com/en/github/automating-your-workflow-with-github-actions/virtual-environments-for-github-actions#creating-and-using-secrets-encrypted-variables](https://help.github.com/en/github/automating-your-workflow-with-github-actions/virtual-environments-for-github-actions#creating-and-using-secrets-encrypted-variables)

### user

Connect as this user.

### connection

Connection type to use.

### timeout

Override the connection timeout in seconds.

### ssh_common_args

Specify common arguments to pass to sftp/scp/ssh.

### sftp_extra_args

Specify extra arguments to pass to sftp only.

### scp_extra_args

Specify extra arguments to pass to scp only.

### ssh_extra_args

Specify extra arguments to pass to ssh only.

### become

Run operations with become.

### become_method

Privilege escalation method to use.

### become_user

Run operations as this user.
required: false

## Example Usage

```yaml
- name: Ansible Playbook
  uses: arillso/action.playbook@master
  with:
    playbook: tests/playbook.yml
    inventory: tests/hosts.yml
    galaxy_file: tests/requirements.yml
  env:
    ANSIBLE_HOST_KEY_CHECKING: 'false'
    ANSIBLE_DEPRECATION_WARNINGS: 'false'
```

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://sbaerlocher.ch"><img src="https://avatars1.githubusercontent.com/u/4160387?v=4" width="100px;" alt=""/><br /><sub><b>Simon Baerlocher</b></sub></a><br /><a href="https://github.com/arillso/action.playbook/commits?author=sbaerlocher" title="Code">💻</a> <a href="#ideas-sbaerlocher" title="Ideas, Planning, & Feedback">🤔</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## License

<!-- markdownlint-disable -->

This project is under the MIT License. See the [LICENSE](licence) file for the full license text.

<!-- markdownlint-enable -->

## Copyright

(c) 2020, Arillso
