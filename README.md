# CLI to generate projects using GoGo library

Nice library is one thing, but what matters is the ability to set up new
projects quickly. If you have just one evening for a new idea, you shouldn't
waste 90% of the time to configure webpack and decide on the deployment
strategy.


## Prerequisites

`gogo` cli won't let you generate a new project if you don't have all the tools
installed. Current requirements are:

- `sqlmigrate` - `go install github.com/rubenv/sql-migrate/...@latest`
- `sqlboiler@v4`:

  ```
  go install github.com/volatiletech/sqlboiler/v4@latest
  go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
  ```

- `watchexec` [site](https://github.com/watchexec/watchexec)
- `flyctl` [docs](https://fly.io/docs/hands-on/install-flyctl/)

Ignore flyctl if you want to roll your own deploy

## Installation

### Install Script

Download `gogo-cli` and install into a local bin directory.

#### MacOS, Linux, WSL

Latest version:

```bash
curl -L https://raw.githubusercontent.com/can3p/gogo-cli/master/generated/install.sh | sh
```

Specific version:

```bash
curl -L https://raw.githubusercontent.com/can3p/gogo-cli/master/generated/install.sh | sh -s 0.0.4
```

The script will install the binary into `$HOME/bin` folder by default, you can override this by setting
`$CUSTOM_INSTALL` environment variable

### Manual download

Get the archive that fits your system from the [Releases](https://github.com/can3p/gogo-cli/releases) page and
extract the binary into a folder that is mentioned in your `$PATH` variable.

## Usage

Let's say we want to create a `bunny` project:

```
cd ~/code
gogo-cli generate bunny --email 'your@email.com' --repo 'github.com/can3p/bunny' --testemail 'your@email.com' --out bunny
echo "SESSION_SALT=random_string" >> bunny/cmd/web/.env
echo "SITE_ROOT=http://localhost:8080" >> bunny/cmd/web/.env
echo "DATABASE_URL=<insert your postgres connection string there>" >> bunny/cmd/web/.env
cd bunny
./sqlmigrate.sh up
./generate.sh
cd cmd/web
yarn
yarn watch # in one tab
make watchexec # in another tab
```

### To Replace in the template

* `{{ .ProjectName }}`
* `{{ .ProjectEmail }}`
* `{{ .ProjectRepo }}`
* `{{ .TestemailHead }}` `john@mail.wat` is split into `john` and `mail.wat`
* `{{ .TestemailTail }}`
