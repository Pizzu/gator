package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Pizzu/gator/internal/database"
	"github.com/google/uuid"
)

// User Handlers
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

	s.logger.Info("User switched successfully!")
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

	s.logger.Info(fmt.Sprintf("User created and set successfully: %v", user))
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
			s.logger.Printf("* %s (current)", user.Name)
		} else {
			s.logger.Printf("%s", user.Name)
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

	s.logger.Info("users deleted successfully")
	return nil
}

// Feed Handlers
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

// Feed Follow handler
func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctx := context.Background()

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

	s.logger.Info(fmt.Sprintf("%s started following %s feed", feedFollow.UserName, feedFollow.FeedName))

	return nil
}

func handlerFeedUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctx := context.Background()

	url := cmd.Args[0]

	err := s.db.UnfollowFeed(ctx, database.UnfollowFeedParams{UserID: user.ID, Url: url})

	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("%s stopped following %s feed", user.Name, url))

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feedsFollowed, err := s.db.GetFeedFollowsForUser(ctx, user.ID)

	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("%s is following:", user.Name))
	for _, feedFollowed := range feedsFollowed {
		s.logger.Printf("- %s", feedFollowed.FeedName)
	}

	return nil
}

// Handler aggregator
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

// Browse Handler
func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Found %d posts for user %s:\n", len(posts), user.Name))
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("Desc: %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
