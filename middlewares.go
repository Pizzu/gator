package main

import (
	"context"
	"errors"

	"github.com/Pizzu/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		currentUsername := s.cfg.CurrentUserName
		if currentUsername == "" {
			return errors.New("not logged in, sign in first")
		}
		ctx := context.Background()
		user, err := s.db.GetUserByName(ctx, currentUsername)

		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
