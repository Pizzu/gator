package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Pizzu/gator/internal/api"
	"github.com/Pizzu/gator/internal/config"
	"github.com/Pizzu/gator/internal/database"
	"github.com/charmbracelet/log"
)

type state struct {
	cfg    *config.Config
	db     *database.Queries
	client api.Client
	logger *log.Logger
}

func NewState(cfg *config.Config, db *database.Queries, logger *log.Logger) *state {
	return &state{
		cfg:    cfg,
		db:     db,
		client: api.NewClient(5 * time.Second),
		logger: logger,
	}
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func Execute(s *state) error {
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetUsers)
	cmds.register("users", handlerGetAllUsers)
	cmds.register("agg", handlerAggregator)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetAllFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFeedFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerFeedUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		return fmt.Errorf("usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	return cmds.run(s, command{Name: cmdName, Args: cmdArgs})
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
