# Command Line Arguments for `gh-digest`

`gh-digest` offers a variety of command line arguments to tailor your GitHub digest experience. Here's a breakdown of what you can use:

### Global Options

- `--env-var`: Specify the name of the environment variable where your GitHub token is stored. Default is `GITHUB_TOKEN`.
- `--repo, -r`: The name of the repository to fetch pull requests from. Use `all` for all repositories or specify a specific repository name.
- `--status, -s`: Filter pull requests based on their status. Options are `open`, `closed`, or `all`.
- `--with-orgs, -w`: Set this flag to fetch pull requests from organizations you're part of.
- `--org, -o`: Specify the name of the organization from which to fetch pull requests.

### Usage

To use these arguments, simply append them to your `gh-digest` command like so:

```bash
gh-digest --repo my-repo --status open --with-orgs
```