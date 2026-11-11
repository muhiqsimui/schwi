package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"schwi/internal/bot"
	"schwi/internal/config"
)

func main() {
	log.Println("[SCHWI] Starting...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[CONFIG] %v", err)
	}

	log.Println("[CONFIG] Configuration loaded")

	// Initialize Discord bot
	discordBot, err := bot.New(
		cfg.DiscordToken,
		cfg.DiscordGuildID,
	)
	if err != nil {
		log.Fatalf("[DISCORD] Initialization failed: %v", err)
	}

	// Register event handlers before connecting
	discordBot.RegisterHandlers()

	// Connect to Discord
	if err := discordBot.Start(); err != nil {
		log.Fatalf("[DISCORD] Connection failed: %v", err)
	}

	// Register slash commands
	if err := discordBot.RegisterCommands(); err != nil {
		log.Printf("[DISCORD] Command registration failed: %v", err)

		_ = discordBot.Stop()
		os.Exit(1)
	}

	log.Println("[SCHWI] Bot is ready!")

	// Wait for shutdown signal
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	log.Println("[SCHWI] Shutting down...")

	if err := discordBot.Stop(); err != nil {
		log.Printf("[DISCORD] Shutdown error: %v", err)
	}

	log.Println("[SCHWI] Shutdown complete")
}
