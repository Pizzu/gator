package main

import (
	"context"
	"fmt"
	"time"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)

	if err != nil {
		return err
	}

	markedFeed, err := s.db.MarkFeedFetched(ctx, feed.ID)

	if err != nil {
		return err
	}

	fetchedFeed, err := s.client.FetchFeed(ctx, markedFeed.Url)

	if err != nil {
		return err
	}

	for _, post := range fetchedFeed.Channel.Item {
		fmt.Printf("Post: %+v\n", post.Title)
	}

	return nil
}
