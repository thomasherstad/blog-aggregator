package main

import (
	"context"
	"fmt"

	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("logged in user is not in database, error: %w", err)
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
