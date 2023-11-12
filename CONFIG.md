### `config.md` - Config File Details

```markdown
# Configuration File for `gh-digest`

`gh-digest` uses a `config.toml` file for additional configurations. Below are the sections and settings you can adjust:

### `[print]`

- `style`: Specifies the print style for your output. Available ones at: `./internal/types.go`


### `[rate]`

- `type`: Can be either `percentage` or `fixed`. This determines the rate limit type for GitHub API requests.
- `value`: The value of the rate limit. If `type` is `percentage`, this is a percentage; if `type` is `fixed`, this is a numeric limit.

### Example `config.toml`

```toml
title = "GH Daily Digest CLI"

[print]
style = ""

[rate]
type = "percentage"
value = 10
