package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/rss+xml")
	req.Header.Set("User-Agent", "gator")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res == nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	defer res.Body.Close()

	rawData, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed

	err = xml.Unmarshal(rawData, &rssFeed)

	if err != nil {
		return nil, err
	}

	// Decode escaped HTML entities
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}
