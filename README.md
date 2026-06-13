# ghir (GitHub Immutable Releases)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/suzuki-shunsuke/ghir)
[Install](INSTALL.md)

ghir is a CLI making past GitHub Releases immutable.

About GitHub Immutable Releases, please see the following links:

- https://github.blog/changelog/2025-08-26-releases-now-support-immutability-in-public-preview/
- https://github.com/orgs/community/discussions/171210

Immutable Releases protect your software supply chain by preventing any changes to released assets.
While enabling Immutable Releases is straightforward, previously created releases remain vulnerable.
ghir is a CLI tool that secures your past releases by making them immutable.

## How To Use

0. Enable Immutable Releases
1. Run ghir

```sh
ghir [--log-level <debug|info|warn|error>] <repo full name>
```

e.g.

```sh
ghir aquaproj/aqua
```

## GitHub Access Token

ghir requires a GitHub Access Token.

- Required Permissions: `contents:write`
- Scopes (accessible repositories): A repository to be updated

The priority of the access token is as follows:

1. Environment Variables
    1. `GHIR_GITHUB_TOKEN`
    1. `GITHUB_TOKEN`
1. [ghtkn integration](https://github.com/suzuki-shunsuke/ghtkn)
    1. `export GHIR_ENABLE_GHTKN=true`
    1. `export GHTKN_ENABLE=true`
    1. If the configuration file of ghtkn exists, the integration is enabled

## How It Works

1. Get GitHub Releases by GitHub API
1. Exclude draft releases and immutable releases
1. Update releases without any parameters by GitHub API to make all releases immutable

## ProTip: Run ghir for multiple repositories

ghir alone can only be executed in a single repository.

However, by combining other tools, you can run ghir against multiple repositories.

### Example 1: Use repository list file

```sh
cat repos.txt
username_or_orgname/foo
username_or_orgname/bar

cat repos.txt |
  xargs -n 1 ghir
```

### Example 2: Use `gh repo list`

[gh repo list](https://cli.github.com/manual/gh_repo_list)

```sh
gh repo list [<owner>] \
  --source \
  --no-archived \
  --json nameWithOwner \
  --template '{{range .}}{{.nameWithOwner}}{{"\n"}}{{end}}' \
  --limit 100 |
  xargs -n 1 ghir
```
