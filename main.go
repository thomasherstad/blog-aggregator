package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
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

	currentState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	var commands commands
	commands.actions = make(map[string]func(*state, command) error)
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("users", handlerUsers)
	commands.register("reset", handlerReset)

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

	err = commands.run(&currentState, currentCommand)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Current username: %s\n", cfg.CurrentUsername)
	fmt.Println(*currentState.cfg)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("a username is required")
	}

	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), sql.NullString{String: username, Valid: true})
	if err != nil {
		return errors.New("the username does not exist in the database")
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("User has been set to:", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("a username is required")
	}
	username := cmd.args[0]

	dbUser, err := s.db.GetUser(context.Background(), sql.NullString{String: username, Valid: true})
	if err == nil {
		fmt.Println("Username already found in database:", dbUser)
		return errors.New("the username already exists")
	}

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: sql.NullString{
			String: username,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("New user successfully stored in database:", newUser)

	handlerLogin(s, cmd)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("the reset command does not take any arguments")
	}

	s.db.DeleteUsers(context.Background())
	fmt.Println("All users deleted")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("the users command does not take any arguments")
	}

	usersList, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(usersList) == 0 {
		fmt.Println("There are currently no registered users.")
		return nil
	}

	for _, user := range usersList {
		fmt.Printf("- %s ", user.Name.String)
		if user.Name.String == s.cfg.CurrentUsername {
			fmt.Printf("(current)")
		}
		fmt.Printf("\n")
	}

	return nil
}
