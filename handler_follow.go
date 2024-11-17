package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow takes 1 argument - %v given", len(cmd.args))
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
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("problem creating feed follow, error: %w", err)
	}

	fmt.Printf("%s is now following %s\n", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following takes no arguments - %v given", len(cmd.args))
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("problem getting feeds, error: %w", err)
	}

	for _, feed := range followedFeeds {
		fmt.Println("-", feed.FeedName)
	}

	return nil
}
