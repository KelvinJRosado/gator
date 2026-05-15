package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/KelvinJRosado/gator/internal/command"
	"github.com/KelvinJRosado/gator/internal/config"
	"github.com/KelvinJRosado/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// Load config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error opening config: %v", err)
	}

	// Connect to DB
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	// Create state
	st := command.State{
		Db:  dbQueries,
		Cfg: cfg,
	}

	// Create commands registry
	cmds := command.CreateCommandsRegistry()
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)
	cmds.Register("reset", command.HandlerReset)
	cmds.Register("users", command.HandlerUsers)
	cmds.Register("agg", command.HandlerAgg)
	cmds.Register("addfeed", command.HandlerAddFeed)

	// Grab CLI args
	if len(os.Args) < 2 {
		log.Fatalf("Insufficient args given (min. 2)")
	}
	args := os.Args[1:]

	userCommand := command.Command{
		Name: args[0],
		Args: args[1:],
	}

	err = cmds.Run(&st, userCommand)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

}
