package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ndious/delivr/internal/command"
	"github.com/ndious/delivr/internal/config"
	"github.com/ndious/delivr/internal/discord"
	"github.com/ndious/delivr/internal/logger"
)

func main() {
	// Parse command line flags
	daemonMode := flag.Bool("daemon", false, "Run in daemon mode (don't exit after running commands)")
	configPath := flag.String("config", "", "Path to the configuration file (default: config.json in the current directory or ~/.delivr/config.json)")
	initConfig := flag.Bool("init", false, "Generate a default configuration file")
	outPath := flag.String("out", "config.json", "Path for the generated configuration file when using --init")
	flag.Parse()

	// Check if we should generate a default configuration file
	if *initConfig {
		log.Printf("Generating default configuration file at: %s", *outPath)
		if err := config.CreateDefaultConfig(*outPath); err != nil {
			log.Fatalf("Failed to create default configuration: %v", err)
		}
		log.Printf("Default configuration created successfully. Please edit %s with your Discord credentials.", *outPath)
		return
	}

	// Initialize logger
	log.SetOutput(os.Stdout)
	log.Println("Starting Delivr - Docker Command Runner with Discord Integration")

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded from: %s", config.GetLoadedConfigPath())

	// Initialize Discord client
	discord, err := discord.NewClient(cfg.Discord.Token, cfg.Discord.ChannelID)
	if err != nil {
		log.Fatalf("Failed to initialize Discord client: %v", err)
	}

	// Send startup message
	if err := discord.SendMessage("🚀 Delivr service started"); err != nil {
		log.Printf("Warning: Could not send startup message: %v", err)
	}

	// Initialize logger
	cmdLogger, err := logger.NewCommandLogger(cfg.Logs)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer cmdLogger.Close()

	// Initialize Docker runner with the global working directory and docker host
	dockerHost := ""
	if cfg.Docker.Host != "" {
		dockerHost = cfg.Docker.Host
	}
	cmdRunner := command.NewRunner(discord, cmdLogger, cfg.WorkingDir, dockerHost)

	// Execute commands defined in config
	for _, cmd := range cfg.Commands {
		if err := cmdRunner.Execute(cmd); err != nil {
			log.Printf("Error executing command '%s': %v", cmd.Name, err)
			if err := discord.SendMessage(fmt.Sprintf("❌ Error executing command '%s': %v", cmd.Name, err)); err != nil {
				log.Printf("Failed to send error message to Discord: %v", err)
			}
		}
	}

	// If not in daemon mode, exit after running commands
	if !*daemonMode {
		// Send shutdown message
		if err := discord.SendMessage("✅ Delivr - Toutes les commandes ont été exécutées"); err != nil {
			log.Printf("Warning: Could not send completion message: %v", err)
		}
		log.Println("All commands executed, shutting down...")
		return
	}

	// In daemon mode, setup signal handling for graceful shutdown
	log.Println("Running in daemon mode, press Ctrl+C to exit")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigCh
	log.Printf("Received signal %v, shutting down...", sig)

	// Send shutdown message
	if err := discord.SendMessage("🛑 Delivr service stopping"); err != nil {
		log.Printf("Warning: Could not send shutdown message: %v", err)
	}

	log.Println("Shutdown complete")
}
