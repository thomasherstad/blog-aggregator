package main

import (
	"fmt"
	"time"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg takes 1 argument - %v given", len(cmd.args))
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("problem parsing time between requests, error: %w", err)
	}

	fmt.Println("Collecting feeds every", timeBetweenReqs)

	// collect feeds every x time
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}
