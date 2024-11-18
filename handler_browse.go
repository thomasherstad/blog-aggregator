package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/thomasherstad/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	if len(cmd.args) == 1 {
		limit = 2
	} else if len(cmd.args) == 2 {
		var err error
		limit, err = strconv.Atoi(cmd.args[1])
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("browse takes 1 or 2 arguments, %v given", len(cmd.args))
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user, error: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n", post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
