package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("a username is required")
	}

	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
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

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	fmt.Println("New user successfully stored in database:", newUser)

	handlerLogin(s, cmd)
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
		fmt.Printf("- %s ", user.Name)
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("(current)")
		}
		fmt.Printf("\n")
	}

	return nil
}
