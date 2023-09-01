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


## Install

TBD

## Usage

TBD

Something like this, let's say we want to create a `bunny` project:

```
cd ~/code
gogo generate project bunny
# or: gogo generate project bunny --skip-deploy
cd bunny/cmd/web
make watchexec
```

### To Replace in the template

* `<projectname>`
* `<projectemail>`
* `<projectrepo>`
* `<testemailhead>` `john@mail.wat` is split into `john` and `mail.wat`
* `<testemailtail>`
