package main

import (
	"log"
	"os"
)

func main() {
	log.Printf("[*] Starting st3fan/mijia-hub ")

	// If we are running under systemd then no need to log datestamps
	if os.Getenv("INVOCATION_ID") != "" {
		log.SetFlags(0)
	}

	cfg, err := newConfigurationFromEnvironment()
	if err != nil {
		log.Fatalf("[F] Failed to parse config from environment: %s", err)
	}

	app, err := newApplication(cfg)
	if err != nil {
		log.Fatalf("[F] Failed to create application: %v", err)
	}

	if err := app.run(); err != nil {
		log.Fatalf("[F] Failed to run application: %s", err)
	}
}
