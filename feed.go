package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
