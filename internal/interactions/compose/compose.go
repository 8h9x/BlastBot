package compose

import "github.com/disgoorg/disgo/discord"

func Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "compose",
		Description: "Advanced command for power users to have more direct access to API calls.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "raw",
				Description: "Request builder for URLs included in the epic API whitelist.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "mcp",
				Description: "Request builder for fortnite MCP operations.",
			},
		},
	}
}
