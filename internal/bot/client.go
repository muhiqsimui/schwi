package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	Session *discordgo.Session
	GuildID string
}

func New(token, guildID string) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	session.Identify.Intents = discordgo.IntentsGuilds

	client := &Client{
		Session: session,
		GuildID: guildID,
	}

	return client, nil
}

func (c *Client) Start() error {
	if err := c.Session.Open(); err != nil {
		return fmt.Errorf("open discord connection: %w", err)
	}

	log.Println("[DISCORD] Connected successfully")
	return nil
}

func (c *Client) Stop() error {
	if c.Session == nil {
		return nil
	}

	if err := c.Session.Close(); err != nil {
		return fmt.Errorf("close discord connection: %w", err)
	}

	log.Println("[DISCORD] Connection closed")
	return nil
}

func (c *Client) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Check bot response",
		},
		{
			Name:        "status",
			Description: "Show server resource usage",
		},
		{
			Name:        "cpu",
			Description: "Show CPU usage information",
		},
		{
			Name:        "ram",
			Description: "Show RAM usage information",
		},
		{
			Name:        "disk",
			Description: "Show disk usage information",
		},
		{
			Name:        "uptime",
			Description: "Show server uptime",
		},
	}

	for _, command := range commands {
		_, err := c.Session.ApplicationCommandCreate(
			c.Session.State.User.ID,
			c.GuildID,
			command,
		)
		if err != nil {
			return fmt.Errorf(
				"register command %s: %w",
				command.Name,
				err,
			)
		}

		log.Printf("[DISCORD] Registered command: /%s", command.Name)
	}

	return nil
}
