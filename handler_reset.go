package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("the reset command does not take any arguments")
	}

	s.db.DeleteAllUsers(context.Background())
	fmt.Println("All users deleted")
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Feeds:", feeds)

	return nil
}
