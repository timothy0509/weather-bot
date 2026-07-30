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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "show",
					Description: "Show current settings",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "顯示目前設定",
						discordgo.ChineseCN: "显示目前设定",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "language",
					Description: "Set language mode",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "設定語言模式",
						discordgo.ChineseCN: "设定语言模式",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "mode",
							Description: "en, tc, sc, bilingual",
							DescriptionLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseTW: "en, tc, sc, bilingual",
								discordgo.ChineseCN: "en, tc, sc, bilingual",
							},
							Type:     discordgo.ApplicationCommandOptionString,
							Required: true,
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
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "設定或停用提示頻道",
						discordgo.ChineseCN: "设定或停用提示频道",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "Channel for weather alerts (omit to disable)",
							DescriptionLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseTW: "天氣提示頻道（留空以停用）",
								discordgo.ChineseCN: "天气提示频道（留空以停用）",
							},
							Type:         discordgo.ApplicationCommandOptionChannel,
							Required:     false,
							ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						},
					},
				},
				{
					Name:        "tide-station",
					Description: "Set default tide station",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "設定預設潮汐站",
						discordgo.ChineseCN: "设定默认潮汐站",
					},
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "station",
							Description:  "3-letter station code",
							DescriptionLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseTW: "3 字母站點代碼",
								discordgo.ChineseCN: "3 字母站点代码",
							},
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Name:        "status",
					Description: "Toggle bot status activity",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.ChineseTW: "切換機械人狀態活動",
						discordgo.ChineseCN: "切换机器人状态活动",
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "enabled",
							Description: "Enable or disable bot status",
							DescriptionLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseTW: "啟用或停用機械人狀態",
								discordgo.ChineseCN: "启用或停用机器人状态",
							},
							Type:     discordgo.ApplicationCommandOptionBoolean,
							Required: true,
						},
					},
				},
			},
		},
	}

	if guildID != "" {
		_, err := s.ApplicationCommandBulkOverwrite(appID, guildID, commands)
		return err
	}
	_, err := s.ApplicationCommandBulkOverwrite(appID, "", commands)
	return err
}
