package bot

import (
	"github.com/bwmarrin/discordgo"
)

var manageServerPermission int64 = discordgo.PermissionManageServer

func registerCommands(s *discordgo.Session, appID, guildID string) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "weather",
			Description: "Get Hong Kong weather information",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "current",
					Description: "Current weather conditions",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "forecast",
					Description: "3-day weather forecast",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "forecast9",
					Description: "9-day weather forecast",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "warnings",
					Description: "Active weather warnings summary",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "warning",
					Description: "Detailed warning information",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "type",
							Description: "Warning code, e.g. WTCSGNL, WRAIN, WTS",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "rain",
					Description: "Hourly rainfall in the past hour",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "uv",
					Description: "Current UV index",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "tide",
					Description: "Hourly astronomical tide heights",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "station",
							Description: "Tide station code (default QUB)",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "lunar",
					Description: "Today's lunar calendar",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "earthquake",
					Description: "Recent earthquake report",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:                     "setup",
			Description:              "Configure the bot for this server",
			DefaultMemberPermissions: &manageServerPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "show",
					Description: "Show current settings",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "language",
					Description: "Set language mode",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "mode",
							Description: "en, tc, sc, bilingual",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{Name: "English", Value: "en"},
								{Name: "Traditional Chinese", Value: "tc"},
								{Name: "Simplified Chinese", Value: "sc"},
								{Name: "Bilingual", Value: "bilingual"},
							},
						},
					},
				},
				{
					Name:        "alerts",
					Description: "Set or disable alerts channel",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "Channel for weather alerts (omit to disable)",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    false,
							ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						},
					},
				},
				{
					Name:        "tide-station",
					Description: "Set default tide station",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "station",
							Description: "3-letter station code",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "status",
					Description: "Toggle bot status activity",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "enabled",
							Description: "Enable or disable bot status",
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Required:    true,
						},
					},
				},
			},
		},
	}

	for _, cmd := range commands {
		var err error
		if guildID != "" {
			_, err = s.ApplicationCommandCreate(appID, guildID, cmd)
		} else {
			_, err = s.ApplicationCommandCreate(appID, "", cmd)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
