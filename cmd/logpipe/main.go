package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/logpipe/internal/config"
)

const defaultConfigPath = "logpipe.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "path to logpipe YAML config file")
	validateOnly := flag.Bool("validate", false, "validate config and exit")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	config.ApplyDefaults(cfg)

	if *validateOnly {
		fmt.Println("config is valid")
		os.Exit(0)
	}

	log.Printf("logpipe starting — listening on %s:%d", cfg.Server.Addr, cfg.Server.Port)
	log.Printf("sources: %d, sinks: %d", len(cfg.Sources), len(cfg.Sinks))

	// TODO: initialise source readers, filter pipeline, and sink writers.
	select {}
}
