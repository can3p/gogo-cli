# CLI to generate projects using GoGo library

Nice library is one thing, but what matters is the ability to set up new
projects quickly. If you have just one evening for a new idea, you shouldn't
waste 90% of the time to configure webpack and decide on the deployment
strategy.


## Prerequisites

`gogo` cli won't let you generate a new project if you don't have all the tools
installed. Current requirements are:

- **Docker** - Required for `db generate` command (runs sqlboiler in a container)
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

Let's say we want to create a [`gogo-cli-example`](https://github.com/can3p/gogo-cli-example) project:

```
cd ~/code
gogo-cli generate gogo-cli-example --email 'your@email.com' --repo 'github.com/can3p/gogo-cli-example' --testemail 'your@email.com' --out gogo-cli-example
echo "SESSION_SALT=random_string" >> gogo-cli-example/cmd/web/.env
echo "SITE_ROOT=http://localhost:8080" >> gogo-cli-example/cmd/web/.env
echo "DATABASE_URL=<insert your postgres connection string there>" >> gogo-cli-example/cmd/web/.env
cd gogo-cli-example
./sqlmigrate.sh up
gogo-cli db generate
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

## Database Commands

### `gogo-cli db generate`

Generates SQLBoiler ORM code from your database schema. This command:

1. Starts a PostgreSQL container using Docker
2. Applies all migrations from a clean state
3. Runs sqlboiler to generate Go code
4. Cleans up the container

**Requirements:**
- Docker must be running
- A `gogo.toml` configuration file in your project root

**Configuration (`gogo.toml`):**

```toml
[database]
migrations = "migrations"

[database.sqlboiler]
pkgname = "core"
output = "pkg/model/core"
wipe = true
blacklist = ["migrations", "http_sessions"]
```

**Usage:**

```bash
# From project root
gogo-cli db generate

# Or specify a different directory
gogo-cli db generate --dir /path/to/project
```
