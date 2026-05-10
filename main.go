package main

import (
	"log"

	"github.com/KelvinJRosado/gator/internal/config"
)

func main() {
	// Load config file
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error opening config: %v", err)
	}

	// Update username in config file
	err = c.SetUser("kelvin")
	if err != nil {
		log.Fatalf("error updating username: %v", err)
	}

	// Print config file contents after reading again
	c, err = config.Read()
	if err != nil {
		log.Fatalf("error opening config: %v", err)
	}
	log.Printf("Config: %v", *c)
}
