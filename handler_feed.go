package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed takes 2 arguments - %v given", len(cmd.args))
	}

	name, url := cmd.args[0], cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed in database. error: %w", err)
	}

	fmt.Println("Feed added to database: ", feed)

	// Create entry in feed_follow db
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("problem following feed, error: %w", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds doesn't take any arguments - %v given", len(cmd.args))
	}

	//get feeds
	rawFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds from database. error: %w", err)
	}

	for _, feed := range rawFeeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get username from database. error: %w", err)
		}

		fmt.Printf("%s: %s - Added by %s\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}
