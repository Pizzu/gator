package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Pizzu/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Collecting feeds every %s...", timeBetweenRequests))

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Couldn't fetch next feed: %v", err))
		return
	}

	markedFeed, err := s.db.MarkFeedFetched(ctx, feed.ID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Couldn't save fetched feed: %v", err))
		return
	}

	fetchedFeed, err := s.client.FetchFeed(ctx, markedFeed.Url)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while fetching feed details: %v", err))
		return
	}

	// Save all the posts for the current feed
	for _, post := range fetchedFeed.Channel.Item {
		publishedAt := sql.NullTime{}

		if t, err := time.Parse(time.RFC1123Z, post.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     post.Title,
			Description: sql.NullString{
				String: post.Description,
				Valid:  true,
			},
			Url:         post.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			s.logger.Error(fmt.Sprintf("Couldn't create post: %v", err))
			continue
		}

	}
	s.logger.Info(fmt.Sprintf("Feed %s collected, %v posts found\n", feed.Name, len(fetchedFeed.Channel.Item)))
}
