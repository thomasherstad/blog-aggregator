package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thomasherstad/blog-aggregator/internal/database"
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

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFecth(context.Background())
	if err != nil {
		return fmt.Errorf("problem getting next feed, error: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("problem marking feed as fetched, error: %w", err)
	}

	err = scrapeFeed(s, nextFeed)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func scrapeFeed(s *state, feed database.Feed) error {
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("problem fetching feed, error: %w", err)
	}

	fmt.Printf("-------Fetched %s, %v posts found-------\n", rssFeed.Channel.Title, len(rssFeed.Channel.Item))
	for _, item := range rssFeed.Channel.Item {
		publishedDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}

		fmt.Printf(item.Title)
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishedDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Printf("Couldn't save %s, error: %s\n", item.Title, err)
		} else {
			fmt.Printf(" - saved\n")
		}
	}
	fmt.Println("-------------------------------")

	return nil

}
