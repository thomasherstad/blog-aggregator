package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Add("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(content, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	unescapeRSSFeed(&feed)
	return &feed, nil
}

func unescapeRSSFeed(feed *RSSFeed) {
	// if feed is not a feed
	// return error

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
}
