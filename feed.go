package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jwoodsiii/blogator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

const timeFormat = "2026-01-29 08:49:12.923622"

func timeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func strToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func timeParse(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func scrapeFeeds(s *state) error {
	// Get the next feed to fetch from the DB.
	// Mark it as fetched.
	// Fetch the feed using the URL (we already wrote this function)
	// Iterate over the items in the feed and print their titles to the console.
	// Update the agg command to now take a single argument: time_between_reqs.

	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch: %v", err)
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)

	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	if err := db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Error fetching feed at url: %v", err)
		return
	}

	// Update your scraper to save posts. Instead of printing out the titles of the posts, save them to the database!
	// If you encounter an error where the post with that URL already exists, just ignore it. That will happen a lot.
	// If it's a different error, you should probably log it.
	// Make sure that you're parsing the "published at" time properly from the feeds.
	// Sometimes they might be in a different format than you expect, so you might need to handle that.
	// You may have to manually convert the data into database/sql types.

	for _, item := range rssFeed.Channel.Item {
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: timeToNullTime(time.Time{}), Title: item.Title, Url: item.Link, Description: strToNullString(item.Description), PublishedAt: timeToNullTime(timeParse(item.PubDate)), FeedID: feed.ID})
		if err != nil {
			if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				log.Printf("Duplicate url, continuing: %v", err)
			} else {
				log.Fatalf("Error creating post: %v", err)
			}
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	if feedUrl == "" {
		return &RSSFeed{}, fmt.Errorf("Need url to fetch")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error generating request: %v", err)
	}

	req.Header.Set("User-Agent", "blogator")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error executing request: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error reading data from request: %v", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("Error unmarshaling xml data: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for idx, i := range feed.Channel.Item {
		feed.Channel.Item[idx].Title = html.UnescapeString(i.Title)
		feed.Channel.Item[idx].Description = html.UnescapeString(i.Description)
	}
	return &feed, nil
}
