package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Pizzu/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFeedFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	currentUsername := s.cfg.CurrentUserName
	if currentUsername == "" {
		return errors.New("not logged in, sign in first")
	}
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, currentUsername)

	if err != nil {
		return err
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(ctx, url)

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

	fmt.Printf("%s started following %s feed\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	currentUsername := s.cfg.CurrentUserName
	if currentUsername == "" {
		return errors.New("not logged in, sign in first")
	}
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, currentUsername)

	if err != nil {
		return err
	}

	feedsFollowed, err := s.db.GetFeedFollowsForUser(ctx, user.ID)

	if err != nil {
		return err
	}

	fmt.Printf("%s is following:\n", user.Name)
	for _, feedFollowed := range feedsFollowed {
		fmt.Printf("%s\n", feedFollowed.FeedName)
	}

	return nil
}
