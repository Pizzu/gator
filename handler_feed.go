package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Pizzu/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	ctx := context.Background()

	feed, err := s.client.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")

	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
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

	fmt.Printf("%+v\n", feed)

	return nil
}
