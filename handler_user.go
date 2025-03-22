package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizzu/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, name)

	if err != nil {
		return fmt.Errorf("couldn't login user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	ctx := context.Background()

	userPayload := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name}

	user, err := s.db.CreateUser(ctx, userPayload)

	if err != nil {
		return fmt.Errorf("couldn't register user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User created and set successfully: %v\n", user)
	return nil
}

func handlerGetAllUsers(s *state, _ command) error {
	ctx := context.Background()

	users, err := s.db.GetUsers(ctx)

	if err != nil {
		return fmt.Errorf("couldn't get all users")
	}

	currentUser := s.cfg.CurrentUserName

	for _, user := range users {
		if currentUser == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerResetUsers(s *state, _ command) error {
	ctx := context.Background()

	err := s.db.DeleteAllUsers(ctx)

	if err != nil {
		return fmt.Errorf("couldn't delete users")
	}

	fmt.Println("users deleted successfully")
	return nil
}
