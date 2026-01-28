package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwoodsiii/blogator/internal/database"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error pulling feeds from db: %v", err)
	}

	for _, f := range feeds {
		fmt.Printf("Name: %s\n", f.Name)
		fmt.Printf("Url: %s\n", f.Url)
		user, err := s.db.GetUserFromId(ctx, f.UserID)
		if err != nil {
			return fmt.Errorf("Error pulling user from id")
		}
		fmt.Printf("User: %s\n", user.Name)
	}

	return nil
}

func handlerFollowing(s *state, cmd command, currUser database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), currUser.Name)
	if err != nil {
		return fmt.Errorf("Error pulling user: %s's follows from db: %v", currUser.Name, err)
	}
	for _, f := range follows {
		fmt.Printf("Name: %s\n", f.FeedName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, currUser database.User) error {
	// It takes a single url argument and creates a new feed follow record for the current user.
	// It should print the name of the feed and the current user once the record is created (which the query we just made should support).
	// You'll need a query to look up feeds by URL.
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("Error getting feed: %v", err)
	}

	ff, err := s.db.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: currUser.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %v", err)
	}
	fmt.Printf("Feed name: %s\n", ff.FeedName)
	fmt.Printf("Username: %s\n", ff.UserName)
	return nil
}

func handlerAddFeed(s *state, cmd command, currUser database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}
	ctx := context.Background()
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	//create feed and associate with user
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: feedName, Url: feedUrl, UserID: currUser.ID})
	if err != nil {
		return fmt.Errorf("Error creating feed: %v", err)
	}

	// create feed follow for current user
	_, err = s.db.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{ID: uuid.New(), CreatedAt: feed.CreatedAt, UpdatedAt: feed.UpdatedAt, UserID: feed.UserID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %s", err)
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
