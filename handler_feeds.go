package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jwoodsiii/blogator/internal/database"
)

func handlerBrowse(s *state, cmd command, currUser database.User) error {
	var input string
	//fmt.Printf("args: %v", cmd.Args)
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s (optional) <limit>", cmd.Name)
	}
	if len(cmd.Args) == 1 {
		input = cmd.Args[0]
	} else {
		input = "2"
	}
	lim, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("Error handling limit: %v", err)
	}
	limit := int32(lim)

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{Name: currUser.Name, Limit: limit})
	if err != nil {
		return fmt.Errorf("Error getting posts for user: %v", err)
	}

	for _, item := range posts {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("CreatedAt: %v\n", item.CreatedAt)
		fmt.Printf("Url: %s\n", item.Url)
		fmt.Printf("Description: %v\n", item.Description)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, currUser database.User) error {
	// accept feed url as arg and unfollow it for current user
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("Error: %v attempting to pull feed from db using url: %s", err, url)
	}

	if err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{UserID: currUser.ID, FeedID: feed.ID}); err != nil {
		return fmt.Errorf("Error attempting to delete feed follow: %v", err)
	}
	return nil
}

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
