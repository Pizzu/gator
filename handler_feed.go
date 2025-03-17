package main

import (
	"context"
	"fmt"
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
