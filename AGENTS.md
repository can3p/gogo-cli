# Agent Context for Gogo CLI

## Overview

CLI tool for scaffolding and managing gogo-based web applications.

## Commands

### `gogo-cli generate <projectname>`
Scaffolds a new gogo project from template.

**Flags:**
- `--email` - From address for emails
- `--repo` - Project repository URL (e.g., `github.com/user/project`)
- `--testemail` - Test email address
- `--out` - Output directory

### `gogo-cli db generate`
Generates SQLBoiler ORM code from database schema.

**Process:**
1. Reads `gogo.toml` configuration
2. Starts PostgreSQL container via Docker
3. Applies all migrations from clean state
4. Runs sqlboiler in Docker container (`ghcr.io/can3p/gogo-docker-deps`)
5. Generates Go code to configured output directory

**Flags:**
- `--dir`, `-d` - Project directory (default: current directory)

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

## Project Structure

```
cmd/
  db.go           - db command group
  db_generate.go  - db generate subcommand
  generate.go     - project scaffolding command
  root.go         - root command setup
pkg/
  config/         - gogo.toml configuration parsing
  dbgen/          - SQLBoiler code generation logic
template/         - Project template files for scaffolding
generated/        - Build info and install scripts
```

## Key Packages

### `pkg/config`
Parses `gogo.toml` configuration file.

**Types:**
- `Config` - Root configuration
- `DatabaseConfig` - Database settings including migrations path
- `SQLBoilerConfig` - SQLBoiler-specific settings (pkgname, output, blacklist)

### `pkg/dbgen`
Handles SQLBoiler code generation.

**Types:**
- `Generator` - Orchestrates the generation process

**Key function:**
- `NewGenerator(projectDir, config)` - Creates generator instance
- `Generator.Run()` - Executes full generation workflow

## Dependencies

- `github.com/can3p/gogo` - Core gogo library (testcontainers)
- `github.com/spf13/cobra` - CLI framework
- `github.com/pelletier/go-toml/v2` - TOML parsing

## Docker Image

The `db generate` command uses `ghcr.io/can3p/gogo-docker-deps` which contains:
- sqlboiler
- sqlboiler-psql driver
- sql-migrate

Source: https://github.com/can3p/gogo-docker-deps
