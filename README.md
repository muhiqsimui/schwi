# Schwi

**Schwi** is a lightweight Discord bot for monitoring Linux server resources directly from your Discord server.

Built with Go, Schwi provides real-time server information through Discord slash commands, allowing you to monitor CPU, memory, disk usage, and system uptime without accessing your server terminal.

## Features

- **CPU Monitoring** — View CPU usage, processor model, and logical CPU count.
- **RAM Monitoring** — Monitor total, used, and available memory.
- **Disk Monitoring** — Check disk capacity, used space, and available storage.
- **Server Status** — Get a complete overview of server resources.
- **Uptime Monitoring** — Check how long your server has been running.
- **Discord Slash Commands** — Interact with your server through native Discord commands.
- **Lightweight** — Built with Go for efficient resource consumption.

## Tech Stack

| Technology | Purpose                            |
| ---------- | ---------------------------------- |
| Go         | Main programming language          |
| DiscordGo  | Discord Gateway and slash commands |
| Fiber v3   | HTTP API framework (planned)       |
| gopsutil   | System resource monitoring         |
| godotenv   | Environment configuration          |

## Available Commands

| Command   | Description                    |
| --------- | ------------------------------ |
| `/ping`   | Check bot response             |
| `/status` | Display complete server status |
| `/cpu`    | Display CPU information        |
| `/ram`    | Display memory usage           |
| `/disk`   | Display disk usage             |
| `/uptime` | Display server uptime          |

## Project Structure

```text
schwi/
├── cmd/
│   └── schwi/
│       └── main.go
├── internal/
│   ├── api/
│   ├── bot/
│   │   ├── client.go
│   │   └── commands.go
│   ├── config/
│   │   └── config.go
│   └── monitor/
│       ├── cpu.go
│       ├── memory.go
│       ├── disk.go
│       └── host.go
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Requirements

- Go 1.24 or later
- A Discord account
- A Discord application and bot
- A Discord server where you have permission to add the bot

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/muhiqsimui/schwi.git
cd schwi
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Configure environment variables

Create your `.env` file:

```bash
cp .env.example .env
```

Edit the configuration:

```bash
nano .env
```

Add the following variables:

```env
DISCORD_TOKEN=YOUR_DISCORD_BOT_TOKEN
DISCORD_GUILD_ID=YOUR_DISCORD_SERVER_ID
PORT=3005
```

| Variable           | Description                           |
| ------------------ | ------------------------------------- |
| `DISCORD_TOKEN`    | Discord bot token                     |
| `DISCORD_GUILD_ID` | Discord server (Guild) ID             |
| `PORT`             | HTTP server port reserved for the API |

**Important:** Never share your Discord bot token or commit your `.env` file to a public repository.

### 4. Create a Discord bot

1. Open the [Discord Developer Portal](https://discord.com/developers/applications).
2. Create a new application.
3. Open the **Bot** section.
4. Create your bot and copy its token.
5. Open **OAuth2 → URL Generator**.
6. Select the `bot` and `applications.commands` scopes.
7. Grant the necessary permissions.
8. Open the generated URL and invite your bot to your server.

To obtain your Discord server ID:

1. Open Discord settings.
2. Enable Developer Mode under Advanced.
3. Right-click your server icon.
4. Select **Copy Server ID**.

### 5. Run Schwi

Run directly:

```bash
go run ./cmd/schwi
```

Or build the binary:

```bash
go build -o schwi ./cmd/schwi
```

Run the binary:

```bash
./schwi
```

Once connected, Schwi will register its slash commands in the configured Discord server.

## Development

Format the Go source code:

```bash
gofmt -w .
```

Check the project:

```bash
go vet ./...
```

Build:

```bash
go build -o schwi ./cmd/schwi
```

## Docker Deployment(Optional)

```
docker compose up -d --build
```

## Security

- Keep your Discord bot token private.
- Store credentials in environment variables.
- Never commit `.env` to Git.
- Use `.env.example` as a configuration template.
- Grant only the Discord permissions required by the bot.
