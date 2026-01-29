package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <scrape duration>", cmd.Name)
	}
	timeBetweenReqs := cmd.Args[0]
	dur, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("Error parsing duration: %v", err)
	}
	fmt.Printf("Collecting feeds every %v\n", dur)
	ticker := time.NewTicker(dur)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
