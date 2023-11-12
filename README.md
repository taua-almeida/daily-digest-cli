# 🚀 Welcome to `gh-digest` - Your Daily GitHub Activity Tracker!

Hey there! Tired of juggling between countless tabs to keep track of your GitHub pull requests? Say no more! `gh-digest` to the rescue! This tool is your personal GitHub, looking through repositories to bring you a neatly compiled daily digest of your GitHub activities. Stay updated with your pending reviews, PR statuses, and more, with just a few clicks (or typing - I expect you to do it)

## 🧰 Built with Go and Some Awesome Packages

`gh-digest` is crafted using the power of Go, and it employs some fantastic packages to keep things smooth and efficient:

- Cobra: For that sleek command-line interface.
- go-github/v56: Fetching data from GitHub without this is like trying to eat soup with a fork.
- BurntSushi/toml: Because who doesn't love a well-structured config file? And why toml and not json, or yml, or txt? Well why not?


## 📘 How to Use gh-digest

**1. Installation:**
- Make sure you have Go installed. If not, visit Go's official site for installation instructions.
- Clone this repo and navigate into it.

**2. Setting Up:**
- Create a `config.toml` file as per your preferences. There's a sample in the repo to get you started.
- Set your GitHub token as an environment variable.

For configuration file details, check out the [Configuration Guide](./CONFIG.md).

**3. Running the Tool:**
- Simply run `go run main.go` and voila! Your GitHub digest is ready!
Use command flags like `--repo`, `--status`, and `--with-orgs` to customize your digest. Mix and match to find your perfect GitHub digest!

For more details on command line arguments, see [Command Line Arguments](./ARGS.md).

## Understanding Output:
The tool will output a list of pull requests based on your specified conditions (author, reviewer, assignee).

It also decorates your PRs with additional info like mergeability and CI/CD status. Neat, right?
