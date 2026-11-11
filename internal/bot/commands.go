package bot

import (
	"fmt"
	"log"
	"time"

	"schwi/internal/monitor"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) RegisterHandlers() {
	c.Session.AddHandler(c.handleInteraction)
}

func (c *Client) handleInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		c.handlePing(s, i)
	case "status":
		c.handleStatus(s, i)
	case "cpu":
		c.handleCPU(s, i)
	case "ram":
		c.handleRAM(s, i)
	case "disk":
		c.handleDisk(s, i)
	case "uptime":
		c.handleUptime(s, i)
	}
}

// ========================================
// PING COMMAND
// ========================================

func (c *Client) handlePing(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong! 🏓",
			},
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Ping response error: %v", err)
	}
}

// ========================================
// STATUS COMMAND
// ========================================

func (c *Client) handleStatus(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Status defer error: %v", err)
		return
	}

	cpuStats, cpuErr := monitor.GetCPU()
	memoryStats, memoryErr := monitor.GetMemory()
	diskStats, diskErr := monitor.GetDisk("/")
	hostStats, hostErr := monitor.GetHost()

	var errorMessages string

	if cpuErr != nil {
		errorMessages += "CPU: " + cpuErr.Error() + "\n"
	}

	if memoryErr != nil {
		errorMessages += "Memory: " + memoryErr.Error() + "\n"
	}

	if diskErr != nil {
		errorMessages += "Disk: " + diskErr.Error() + "\n"
	}

	if hostErr != nil {
		errorMessages += "Host: " + hostErr.Error() + "\n"
	}

	if errorMessages != "" {
		_, err = s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Content: &errorMessages,
			},
		)

		if err != nil {
			log.Printf("[COMMAND] Status error response: %v", err)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Schwi Server Monitor",
		Description: fmt.Sprintf(
			"Server: **%s**\nPlatform: **%s %s**",
			hostStats.Hostname,
			hostStats.Platform,
			hostStats.PlatformVersion,
		),
		Color: 0x2ECC71,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "CPU",
				Value: fmt.Sprintf(
					"%.2f%%\n%s\n%d logical CPUs",
					cpuStats.UsagePercent,
					cpuStats.Model,
					cpuStats.Cores,
				),
				Inline: false,
			},
			{
				Name: "Memory",
				Value: fmt.Sprintf(
					"Used: %s (%.2f%%)\nAvailable: %s\nTotal: %s",
					formatBytes(memoryStats.Used),
					memoryStats.UsedPercent,
					formatBytes(memoryStats.Available),
					formatBytes(memoryStats.Total),
				),
				Inline: true,
			},
			{
				Name: "Disk (/)",
				Value: fmt.Sprintf(
					"Used: %s (%.2f%%)\nFree: %s\nTotal: %s",
					formatBytes(diskStats.Used),
					diskStats.UsedPercent,
					formatBytes(diskStats.Free),
					formatBytes(diskStats.Total),
				),
				Inline: true,
			},
			{
				Name: "System",
				Value: fmt.Sprintf(
					"OS: %s\nKernel: %s\nUptime: %s",
					hostStats.OS,
					hostStats.KernelVersion,
					formatUptime(hostStats.Uptime),
				),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Schwi Server Monitoring",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Status response error: %v", err)
	}
}

// ========================================
// CPU COMMAND
// ========================================

func (c *Client) handleCPU(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)

	if err != nil {
		log.Printf("[COMMAND] CPU defer error: %v", err)
		return
	}

	stats, err := monitor.GetCPU()

	if err != nil {
		message := "Gagal mengambil informasi CPU: " + err.Error()

		_, editErr := s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Content: &message,
			},
		)

		if editErr != nil {
			log.Printf("[COMMAND] CPU error response: %v", editErr)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "CPU Monitor",
		Color: 0x3498DB,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("**%.2f%%**", stats.UsagePercent),
				Inline: true,
			},
			{
				Name:   "Logical CPUs",
				Value:  fmt.Sprintf("%d", stats.Cores),
				Inline: true,
			},
			{
				Name:   "Processor",
				Value:  stats.Model,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Schwi Server Monitoring",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		},
	)

	if err != nil {
		log.Printf("[COMMAND] CPU response error: %v", err)
	}
}

// ========================================
// RAM COMMAND
// ========================================

func (c *Client) handleRAM(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)

	if err != nil {
		log.Printf("[COMMAND] RAM defer error: %v", err)
		return
	}

	stats, err := monitor.GetMemory()

	if err != nil {
		message := "Gagal mengambil informasi RAM: " + err.Error()

		_, editErr := s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Content: &message,
			},
		)

		if editErr != nil {
			log.Printf("[COMMAND] RAM error response: %v", editErr)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "RAM Monitor",
		Description: fmt.Sprintf(
			"Penggunaan memori server saat ini: **%.2f%%**",
			stats.UsedPercent,
		),
		Color: 0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Used",
				Value:  fmt.Sprintf("**%s**", formatBytes(stats.Used)),
				Inline: true,
			},
			{
				Name:   "Available",
				Value:  fmt.Sprintf("**%s**", formatBytes(stats.Available)),
				Inline: true,
			},
			{
				Name:   "Total",
				Value:  fmt.Sprintf("**%s**", formatBytes(stats.Total)),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Schwi Server Monitoring",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		},
	)

	if err != nil {
		log.Printf("[COMMAND] RAM response error: %v", err)
	}
}

// ========================================
// DISK COMMAND
// ========================================

func (c *Client) handleDisk(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Disk defer error: %v", err)
		return
	}

	stats, err := monitor.GetDisk("/")

	if err != nil {
		message := "Gagal mengambil informasi disk: " + err.Error()

		_, editErr := s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Content: &message,
			},
		)

		if editErr != nil {
			log.Printf("[COMMAND] Disk error response: %v", editErr)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Disk Monitor",
		Description: fmt.Sprintf(
			"Penggunaan penyimpanan pada path **%s**: **%.2f%%**",
			stats.Path,
			stats.UsedPercent,
		),
		Color: 0xE67E22,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Used",
				Value:  fmt.Sprintf("**%s**", formatBytes(stats.Used)),
				Inline: true,
			},
			{
				Name:   "Free",
				Value:  fmt.Sprintf("**%s**", formatBytes(stats.Free)),
				Inline: true,
			},
			{
				Name:   "Total",
				Value:  fmt.Sprintf("**%s**", formatBytes(stats.Total)),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Schwi Server Monitoring",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Disk response error: %v", err)
	}
}

// ========================================
// UPTIME COMMAND
// ========================================

func (c *Client) handleUptime(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Uptime defer error: %v", err)
		return
	}

	stats, err := monitor.GetHost()

	if err != nil {
		message := "Gagal mengambil informasi uptime: " + err.Error()

		_, editErr := s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Content: &message,
			},
		)

		if editErr != nil {
			log.Printf("[COMMAND] Uptime error response: %v", editErr)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Server Uptime",
		Description: fmt.Sprintf(
			"Server **%s** telah berjalan selama:",
			stats.Hostname,
		),
		Color: 0x1ABC9C,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Uptime",
				Value:  fmt.Sprintf("**%s**", formatUptime(stats.Uptime)),
				Inline: false,
			},
			{
				Name:   "Operating System",
				Value:  stats.Platform,
				Inline: true,
			},
			{
				Name:   "Kernel",
				Value:  stats.KernelVersion,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Schwi Server Monitoring",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		},
	)

	if err != nil {
		log.Printf("[COMMAND] Uptime response error: %v", err)
	}
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func formatBytes(bytes uint64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf(
		"%.2f %ciB",
		float64(bytes)/float64(div),
		"KMGTPE"[exp],
	)
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	return fmt.Sprintf(
		"%d hari, %d jam, %d menit, %d detik",
		days,
		hours,
		minutes,
		remainingSeconds,
	)
}
