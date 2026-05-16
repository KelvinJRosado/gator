# Gator CLI

Gator is a command-line RSS feed aggregator written in Go. It lets you register users, add and follow RSS feeds, aggregate posts into a PostgreSQL database, and browse collected posts from the terminal.

## Prerequisites

You'll need the following installed to run Gator:

- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)

If you plan to run database migrations from this repo, you'll also need `goose`:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Make sure your Go binary directory is on your `PATH`. For many systems, that's:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Installation

Install the Gator CLI with `go install`:

```sh
go install github.com/KelvinJRosado/gator@latest
```

After installation, you should be able to run:

```sh
gator <command> [args...]
```

## Configuration

Gator expects a config file in your home directory named:

```text
~/.gatorconfig.json
```

Create that file with the following structure:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Update `db_url` to match your local PostgreSQL connection string.

For example, if your PostgreSQL database is named `gator`, you can create it with:

```sh
createdb gator
```

## Database Setup

From the project root, run the database migrations:

```sh
make migrate-up
```

The Makefile reads the database URL from `~/.gatorconfig.json`.

Other useful migration commands:

- `make migrate-status` shows migration status
- `make migrate-down` rolls back the latest migration

## Usage

Run commands with the installed CLI:

```sh
gator <command> [args...]
```

Or, if you're developing locally from the project root, run:

```sh
go run . <command> [args...]
```

## Commands

Register a new user:

```sh
gator register alice
```

Log in as an existing user:

```sh
gator login alice
```

List all users:

```sh
gator users
```

Add a feed and automatically follow it:

```sh
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
```

List all feeds:

```sh
gator feeds
```

Follow an existing feed:

```sh
gator follow "https://blog.boot.dev/index.xml"
```

List feeds followed by the current user:

```sh
gator following
```

Unfollow a feed:

```sh
gator unfollow "https://blog.boot.dev/index.xml"
```

Start aggregating feeds on an interval:

```sh
gator agg 1m
```

Browse collected posts:

```sh
gator browse
```

Browse a specific number of posts:

```sh
gator browse 10
```

Reset user data:

```sh
gator reset
```
