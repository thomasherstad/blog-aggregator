package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments given")
	}

	url := "https://wagslane.dev/index.xml"

	fmt.Println("The url is:", url)

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}
