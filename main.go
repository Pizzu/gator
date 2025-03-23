package main

import (
	"database/sql"
	"os"

	"github.com/Pizzu/gator/internal/cmd"
	"github.com/Pizzu/gator/internal/config"
	"github.com/Pizzu/gator/internal/database"
	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
)

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

	programState := cmd.NewState(&cfg, dbQueries, logger)

	if err := cmd.Execute(programState); err != nil {
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
