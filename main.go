package main

import (
	"log"

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

}
