package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/thomasherstad/blog-aggregator/internal/config"
	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}
	dbQueries := database.New(db)

	currentState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	var commands commands
	commands.actions = make(map[string]func(*state, command) error)
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("users", handlerUsers)
	commands.register("reset", handlerReset)
	commands.register("agg", handlerAggregate)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Error: not enough arguments")
		os.Exit(1)
	}

	var checkedArgs []string
	if len(args[1:]) > 0 {
		checkedArgs = args[1:]
	}
	currentCommand := command{
		name: args[0],
		args: checkedArgs,
	}

	err = commands.run(currentState, currentCommand)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
