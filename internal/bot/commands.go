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
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseTW: "獲取香港天氣資訊",
				discordgo.ChineseCN: "获取香港天气资讯",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "current",
					Description: "Current weather conditions",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "現時天氣狀況",
						discordgo.ChineseCN: "现时天气状况",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "forecast",
					Description: "3-day weather forecast",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "三天天氣預報",
						discordgo.ChineseCN: "三天天气预报",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "forecast9",
					Description: "9-day weather forecast",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "九天天氣預報",
						discordgo.ChineseCN: "九天天气预报",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "warnings",
					Description: "Active weather warnings summary",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "生效天氣警告摘要",
						discordgo.ChineseCN: "生效天气警告摘要",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "warning",
					Description: "Detailed warning information",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "詳細警告資訊",
						discordgo.ChineseCN: "详细警告资讯",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "type",
							Description:  "Warning code, e.g. WTCSGNL, WRAIN, WTS",
							DescriptionLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseTW: "警告代碼，例如 WTCSGNL、WRAIN、WTS",
								discordgo.ChineseCN: "警告代码，例如 WTCSGNL、WRAIN、WTS",
							},
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Name:        "rain",
					Description: "Hourly rainfall in the past hour",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "過去一小時每小時雨量",
						discordgo.ChineseCN: "过去一小时每小时雨量",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "uv",
					Description: "Current UV index",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "現時紫外線指數",
						discordgo.ChineseCN: "现时紫外线指数",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "tide",
					Description: "Hourly astronomical tide heights",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "每小時天文潮汐高度",
						discordgo.ChineseCN: "每小时天文潮汐高度",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "station",
							Description:  "Tide station code (default QUB)",
							DescriptionLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseTW: "潮汐站代碼（預設 QUB）",
								discordgo.ChineseCN: "潮汐站代码（默认 QUB）",
							},
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     false,
							Autocomplete: true,
						},
					},
				},
				{
					Name:        "lunar",
					Description: "Today's lunar calendar",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "今日農曆",
						discordgo.ChineseCN: "今日农历",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "earthquake",
					Description: "Recent earthquake report",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "最近地震報告",
						discordgo.ChineseCN: "最近地震报告",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:                     "setup",
			Description:              "Configure the bot for this server",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.ChineseTW: "設定此伺服器的機械人",
				discordgo.ChineseCN: "设定此服务器的机器人",
			},
			DefaultMemberPermissions: &manageServerPermission,
		},
	}

	if guildID != "" {
		_, err := s.ApplicationCommandBulkOverwrite(appID, guildID, commands)
		return err
	}
	_, err := s.ApplicationCommandBulkOverwrite(appID, "", commands)
	return err
}
