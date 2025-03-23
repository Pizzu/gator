package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizzu/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	ctx := context.Background()

	name := cmd.Args[0]
	url := cmd.Args[1]

	feedPayload := database.CreateFeedParams{
		ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		Name: name, Url: url, UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, feedPayload)

	if err != nil {
		return err
	}

	feedFollowPayload := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		UserID: user.ID, FeedID: feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, feedFollowPayload)

	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("%s started following %s feed", feedFollow.UserName, feedFollow.FeedName))
	s.logger.Info(fmt.Sprintf("%+v", feed))

	return nil
}

func handlerGetAllFeeds(s *state, _ command) error {
	ctx := context.Background()

	feeds, err := s.db.GetAllFeeds(ctx)

	if err != nil {
		return fmt.Errorf("couldn't retrieve feeds")
	}

	for _, feed := range feeds {
		s.logger.Printf("%+v", feed)
	}
	return nil
}
