package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/Pizzu/gator/internal/api"
	"github.com/Pizzu/gator/internal/config"
	"github.com/Pizzu/gator/internal/database"
	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
)

type state struct {
	cfg    *config.Config
	db     *database.Queries
	client api.Client
	logger *log.Logger
}

func main() {
	logger := log.New(os.Stderr)

	cfg, err := config.Read()
	if err != nil {
		logger.Fatalf("error reading config: %v", err)
	}

	db, err := openDB(cfg.DbURL)

	if err != nil {
		logger.Fatal(err.Error())
	}

	defer closeDB(db, logger)

	dbQueries := database.New(db)
	client := api.NewClient(5 * time.Second)

	programState := &state{
		cfg:    &cfg,
		db:     dbQueries,
		client: client,
		logger: logger,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
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
		logger.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func closeDB(db *sql.DB, logger *log.Logger) {
	err := db.Close()
	if err != nil {
		logger.Fatal("Error while closing DB connection")
	}
}
