package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow takes 1 argument - %v given", len(cmd.args))
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("problem getting current user, error: %w", err)
	}

	url := cmd.args[0]

	//Get feed by url
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("problem finding feed in database, error: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("problem creating feed follow, error: %w", err)
	}

	fmt.Printf("%s is now following %s\n", currentUser.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following takes no arguments - %v given", len(cmd.args))
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("problem getting current user, error: %w", err)
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("problem getting feeds, error: %w", err)
	}

	for _, feed := range followedFeeds {
		fmt.Println("-", feed.FeedName)
	}

	return nil
}
