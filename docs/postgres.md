# PostgreSQL

PostgreSQL is a production-ready, open-source database. It's a great choice database for many web applications, and as a back-end engineer, it might be the single most important database to be familiar with.

## How does PostgreSQL work?

Postgres, like most other database technologies, is itself a server. It listens for requests on a port (Postgres' default is `:5432`), and responds to those requests. To interact with Postgres, first you will install the server and start it. Then, you can connect to it using a client like [psql](https://www.postgresql.org/docs/current/app-psql.html#:~:text=psql%20is%20a%20terminal%2Dbased,or%20from%20command%20line%20arguments.) or [PGAdmin](https://www.pgadmin.org/).

## 1. Install

### Mac OS

I recommend using [brew](https://brew.sh/) to install PostgreSQL on Mac OS.

```bash
brew install postgresql
```

### Linux (or WSL)

I recommend using apt to install PostgreSQL on Linux (Ubuntu). Here are the [docs from Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql). The basic steps are:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

## 2. Ensure the installation worked

The `psql` command-line utility is the default client for Postgres. Use it to make sure you're on version 14+ of Postgres:

```bash
psql --version
```

## 3. Start the Postgres server in the background

### Mac OS

```bash
brew services start postgresql
```

### Linux (or WSL)

```bash
sudo service postgresql start
```

## 4. Connect to the server using a client

I'm going to recommend using the PGAdmin client. It's a GUI that makes it easy to interact with Postgres and provides a lot of useful features. If you want to use `psql` on the command line instead, you're welcome to do so.

1. Download [PGAdmin here](https://www.pgadmin.org/).
2. Open PGAdmin and create a new server connection. Here are the connection details you'll need:

### Mac OS (with brew)

* Host: `localhost`
* Port: `5432`
* Username: Your Mac OS username
* Password: *leave this blank*

### Linux (or WSL)

If you're on Linux, you have one more step before you can connect with credentials. On your command line run:

```
sudo passwd postgres
```

Enter a new password *and remember it*. Then restart your shell session.

* Host: `localhost`
* Port: `5432`
* Username: `postgres`
* Password: *the password you created*

## 5. Create a database

A single Postgres server can host multiple databases. In the dropdown menu on the left, open the `Localhost` tab, then right click on "databases" and select "create database".

Name it whatever you like, but you'll need to know the name.

## 6. Query the database

Right click on your database's name in the menu on the left, then select "query tool". You should see a new window open with a text editor. In the text editor, type the following query:

```sql
SELECT version();
```

And click the triangle icon (execute/refresh) to run it. If you see a version number, you're good to go!

*PGAdmin is the "thunder client" of Postgres. It's just a GUI that allows you to run ad-hoc queries against your database.*


# Dependencies

We'll be using a couple of tools to help us out:

* [database/sql](https://pkg.go.dev/database/sql): This is part of Go's standard library. It provides a way to connect to a SQL database, execute queries, and scan the results into Go types.
* [sqlc](https://sqlc.dev/): SQLC is an *amazing* Go program that generates Go code from SQL queries. It's not exactly an [ORM](https://www.freecodecamp.org/news/what-is-an-orm-the-meaning-of-object-relational-mapping-database-tools/), but rather a tool that makes working with raw SQL almost as easy as using an ORM.
* [Goose](https://github.com/pressly/goose): Goose is a database migration tool written in Go. It runs migrations from the same SQL files that SQLC uses, making the pair of tools a perfect fit.

## 1. Install SQLC

SQLC is just a command line tool, it's not a package that we need to import. I recommend [installing](https://docs.sqlc.dev/en/latest/overview/install.html) it using `go install`:

```bash
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

Then run `sqlc version` to make sure it's installed correctly.

## 2. Install Goose

Like SQLC, Goose is just a command line tool. I also recommend [installing](https://github.com/pressly/goose#install) it using `go install`:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Run `goose -version` to make sure it's installed correctly.

## 3. Create the `users` migration

I recommend creating an `sql` directory in the root of your project, and in there creating a `schema` directory.

A "migration" is a SQL file that describes a change to your database schema. For now, we need our first migration to create a `users` table. The simplest format for these files is:

```
number_name.sql
```

For example, I created a file in `sql/schema` called `001_users.sql` with the following contents:

```sql
-- +goose Up
CREATE TABLE ...

-- +goose Down
DROP TABLE users;
```

Write out the `CREATE TABLE` statement in full, I left it blank for you to fill in. A `user` should have 4 fields:

* id: a `UUID` that will serve as the primary key
* created_at: a `TIMESTAMP` that can not be null
* updated_at: a `TIMESTAMP` that can not be null
* name: a string that can not be null

The `-- +goose Up` and `-- +goose Down` comments are required. They tell Goose how to run the migration. An "up" migration moves your database from its old state to a new state. A "down" migration moves your database from its new state back to its old state.

By running all of the "up" migrations on a blank database, you should end up with a database in a ready-to-use state. "Down" migrations are only used when you need to roll back a migration, or if you need to reset a local testing database to a known state.

## 4. Run the migration

`cd` into the `sql/schema` directory and run:

```bash
goose postgres CONN up
```

Where `CONN` is the connection string for your database. sThe format is:

```
protocol://username:password@host:port/database
```

Run your migration! Make sure it works by using PGAdmin to find your newly created `users` table.

## 5. Save your connection string as an environment variable

Add your connection string to your `.env` file. When using it with `goose`, you'll use it in the format we just used. However, here in the `.env` file it needs an additional query string:

```
protocol://username:password@host:port/database?sslmode=disable
```

Your application code needs to know to not try to use SSL locally.

## 6. Configure [SQLC](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)

You'll always run the `sqlc` command from the root of your project. Create a file called `sqlc.yaml` in the root of your project. Here is mine:

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

We're telling SQLC to look in the `sql/schema` directory for our schema structure (which is the same set of files that Goose uses, but sqlc automatically ignores "down" migrations), and in the `sql/queries` directory for queries. We're also telling it to generate Go code in the `internal/database` directory.

## 7. Write a query to create a user

Inside the `sql/queries` directory, create a file called `users.sql`. Here is mine:

```sql
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

`$1`, `$2`, `$3`, and `$4` are parameters that we'll be able to pass into the query in our Go code. The `:one` at the end of the query name tells SQLC that we expect to get back a single row (the created user).

Keep the [SQLC docs](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html) handy, you'll probably need to refer to them again later.

## 8. Generate the Go code

Run `sqlc generate` from the root of your project. It should create a new package of go code in `internal/database`.