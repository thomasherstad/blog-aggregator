# Gator - A Blog Aggregator
CLI tool for aggregating, storing and viewing RSS feeds of your choice.  This is a project part of [boot.dev](https://boot.dev)'s backend track to learn more about Golang and PostgreSQL.

## Prerequisites
- Latest Go toolchain
- Local Postgres database

## Quick Start (Linux)
1. Install gator in your chosen directory by running the command `go install github.com/thomasherstad/blog-aggregator` in your terminal. 
2. Create a new PostgreSQL database and name it `gator`.
3. Create a config file called `.gatorconfig.json` in your home directory and add your database link.
    ```json
        {
            "db_url":"postgres://YourPostgresUsername:YourPostgresPassword@localhost:5432/DatabaseName?sslmode=disable",
        }
    ```
    Replace the values with your database connection string.
4. Run the program with `run  go . register <username>`
5. Add a feed with `run go . addfeed <feed name> <feed URL>`
6. Aggregate posts with `run go . agg 1m` and let it save new posts every minute
7. In a new terminal window, run `run go . browse 5` to view the 5 most recent posts

## Commands
### Register a User
`go run . register <your username>`
- Registers the user in the database and automatically logs in the new user

### Login a Different User
`go run . login <your username>`
- Logs in a different user that is already registered

### List All Users
`go run . users`
- Lists all registered users

### List All Feeds
`go run . feeds`
- Lists all feeds in the database and which user created them

### Add a Feed
`go run . addfeed <name of feed> <feed URL>`
- Adds a feed to the database which is automatically followed by the logged-in user

### Follow a Feed
`go run . follow <feed URL>`
- Follow a feed that is already in the database created by another user

### Unfollow a Feed
`go run . unfollow <feed URL>`
- Unfollow a feed the user is following

### List Followed Feeds
`go run . following`
- Lists all feeds the current user is following

### Aggregate Posts From Your Feeds
`go run . agg <time interval>`
- Starts aggregating a feed every time interval. The interval can be set by for example 1h, 1m or 10s

### Browse Posts
`go run . browse <amount>`
- Shows you the _amount_ most recent posts from all the feeds you are following

### Reset the Complete Database (Danger)

`go run . reset`
- Completely wipes all users, feeds and posts from the database

