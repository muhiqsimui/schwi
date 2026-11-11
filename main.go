
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN belum diatur")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize Discord session
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Gagal membuat Discord session: ", err)
	}

	// Enable required intents
	bot.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent

	// Register message handler
	bot.AddHandler(func(
		s *discordgo.Session,
		m *discordgo.MessageCreate,
	) {
		// Ignore messages from bots
		if m.Author == nil || m.Author.Bot {
			return
		}

		message := strings.ToLower(strings.TrimSpace(m.Content))

		// Handle ping command
		if message == "ping" || message == "!ping" {
			_, err := s.ChannelMessageSend(
				m.ChannelID,
				"Pong! 🏓",
			)

			if err != nil {
				log.Println("Gagal mengirim pesan:", err)
			}
		}
	})

	// Connect to Discord
	if err := bot.Open(); err != nil {
		log.Fatal("Gagal terhubung ke Discord: ", err)
	}
	defer bot.Close()

	log.Println("Discord bot berhasil terhubung!")

	// Initialize Fiber
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     "Go Discord Bot",
			"status":  "running",
			"discord": "connected",
		})
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Start Fiber HTTP server
	httpErr := make(chan error, 1)

	go func() {
		httpErr <- app.Listen(":" + port)
	}()

	log.Printf("Fiber berjalan di http://localhost:%s", port)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("Menghentikan aplikasi...")
	case err := <-httpErr:
		if err != nil {
			log.Println("Fiber berhenti:", err)
		}
	}

	if err := app.Shutdown(); err != nil {
		log.Println("Fiber shutdown error:", err)
	}

	log.Println("Aplikasi berhasil dihentikan")
}
