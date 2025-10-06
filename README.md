# ghir (GitHub Immutable Releases)

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
ghir [--log-level <debug|info|warn|error>] [--enable-ghtkn] <repo full name>
```

e.g.

```sh
ghir aquaproj/aqua
```

## GitHub Access Token

ghir requires a GitHub Access Token.

- Required Permissions: `contents:write`
- Scopes (accessible repositories): A repository to be updated

Environment Variables

1. `GHIR_GITHUB_TOKEN`
1. `GITHUB_TOKEN`

Or you can also use [ghtkn integration](https://github.com/suzuki-shunsuke/ghtkn).

```sh
ghir --enable-ghtkn <repo>
```

Or

```sh
export GHIR_ENABLE_GHTKN=true
```

## How It Works

1. Get GitHub Releases by GitHub API
1. Exclude draft releases and immutable releases
1. Update releases without any parameters by GitHub API to make all releases immutable

## Note

Even if ghir makes releases immutable, release attestations aren't created.
I sent a feature request.
[For more details, please see the discussion.](https://github.com/orgs/community/discussions/171210#discussioncomment-14601356)

## LICENSE

[MIT](LICENSE)
