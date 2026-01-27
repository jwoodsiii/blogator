package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage %s", cmd.Name)
	}

	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()

	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("Error fetching rss feed: %v", err)
	}
	fmt.Printf("Channel Title: %s\n", feed.Channel.Title)
	fmt.Printf("Channel Link: %s\n", feed.Channel.Link)
	fmt.Printf("Channel Description: %s\n", feed.Channel.Description)
	for _, i := range feed.Channel.Item {
		fmt.Printf("Item Title: %s\n", i.Title)
		fmt.Printf("Item Link: %s\n", i.Link)
		fmt.Printf("Item Description: %s\n", i.Description)
		fmt.Printf("Item PubDate: %s\n", i.PubDate)
	}
	return nil
}
