package main

import (
	"log"
	"os"

	"github.com/KelvinJRosado/gator/internal/command"
	"github.com/KelvinJRosado/gator/internal/config"
)

func main() {
	// Load config file
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error opening config: %v", err)
	}

	// Create state
	st := command.State{
		Cfg: c,
	}

	// Create commands registry
	cmds := command.CreateCommandsRegistry()
	cmds.Register("login", command.HandlerLogin)

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
