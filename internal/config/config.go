package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken   string
	DiscordGuildID string
	Port           string
}

func Load() (*Config, error) {
	// Load environment variables from .env
	_ = godotenv.Load()

	cfg := &Config{
		DiscordToken:   strings.TrimSpace(os.Getenv("DISCORD_TOKEN")),
		DiscordGuildID: strings.TrimSpace(os.Getenv("DISCORD_GUILD_ID")),
		Port:           strings.TrimSpace(os.Getenv("PORT")),
	}

	// Validate Discord token
	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN tidak ditemukan")
	}

	// Validate Discord Guild ID
	if cfg.DiscordGuildID == "" {
		return nil, fmt.Errorf("DISCORD_GUILD_ID tidak ditemukan")
	}

	if _, err := strconv.ParseUint(cfg.DiscordGuildID, 10, 64); err != nil {
		return nil, fmt.Errorf("DISCORD_GUILD_ID harus berupa angka: %w", err)
	}

	// Default port
	if cfg.Port == "" {
		cfg.Port = "3005"
	}

	// Validate port
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("PORT tidak valid: %s", cfg.Port)
	}

	return cfg, nil
}
