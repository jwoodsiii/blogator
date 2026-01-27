package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwoodsiii/blogator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}
	ctx := context.Background()
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	currUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error pulling current user from db: %v", err)
	}

	//create feed and associate with user
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: feedName, Url: feedUrl, UserID: currUser.ID})
	if err != nil {
		return fmt.Errorf("Error creating feed: %v", err)
	}

	// print record fields
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("FeedName: %s\n", feed.Name)
	fmt.Printf("FeedUrl: %s\n", feed.Url)
	fmt.Printf("UserId: %s\n", feed.UserID)

	return nil
}

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
