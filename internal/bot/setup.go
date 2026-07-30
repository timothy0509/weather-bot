package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"weather-bot/internal/hko"
	"weather-bot/internal/i18n"
)

func (b *Bot) handleSetupCommand(i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	lang := b.guildLanguage(i.GuildID)

	if len(data.Options) == 0 {
		b.respond(i, i18n.T("unknown_setup_command", lang))
		return
	}
	op := data.Options[0]

	switch op.Name {
	case "show":
		b.handleSetupShow(i, lang)
	case "language":
		b.handleSetupLanguage(i, lang, op)
	case "alerts":
		b.handleSetupAlerts(i, lang, op)
	case "tide-station":
		b.handleSetupTideStation(i, lang, op)
	case "status":
		b.handleSetupStatus(i, lang, op)
	default:
		b.respond(i, i18n.T("unknown_setup_command", lang))
	}
}

func (b *Bot) handleSetupShow(i *discordgo.InteractionCreate, lang i18n.Language) {
	gs, err := b.DB.GetGuildSettings(i.GuildID)
	if err != nil {
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	channelValue := i18n.T("alert_channel_disabled", lang)
	if gs.AlertChannelID.Valid {
		channelValue = fmt.Sprintf("<#%s>", gs.AlertChannelID.String)
	}

	statusValue := i18n.T("disabled", lang)
	if gs.BotStatusEnabled {
		statusValue = i18n.T("enabled", lang)
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("settings", lang),
		Color: 0x95A5A6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   i18n.T("language", lang),
				Value:  gs.Language,
				Inline: true,
			},
			{
				Name:   i18n.T("alert_channel", lang),
				Value:  channelValue,
				Inline: true,
			},
			{
				Name:   i18n.T("station", lang),
				Value:  gs.TideStation,
				Inline: true,
			},
			{
				Name:   i18n.T("bot_status", lang),
				Value:  statusValue,
				Inline: true,
			},
		},
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleSetupLanguage(i *discordgo.InteractionCreate, lang i18n.Language, op *discordgo.ApplicationCommandInteractionDataOption) {
	mode := ""
	for _, opt := range op.Options {
		if opt.Name == "mode" {
			mode = opt.StringValue()
		}
	}
	if !i18n.IsValid(mode) {
		b.respond(i, i18n.T("invalid_language_mode", lang))
		return
	}
	if err := b.DB.SetLanguage(i.GuildID, mode); err != nil {
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}
	b.respond(i, fmt.Sprintf("%s **%s**.", i18n.T("language_set", i18n.Normalize(mode)), mode))
}

func (b *Bot) handleSetupAlerts(i *discordgo.InteractionCreate, lang i18n.Language, op *discordgo.ApplicationCommandInteractionDataOption) {
	if len(op.Options) == 0 {
		if err := b.DB.SetAlertChannel(i.GuildID, ""); err != nil {
			b.respond(i, i18n.T("error_fetching_data", lang))
			return
		}
		b.respond(i, i18n.T("alert_channel_removed", lang))
		return
	}
	for _, opt := range op.Options {
		if opt.Name == "channel" {
			channelID, ok := opt.Value.(string)
			if !ok {
				b.respond(i, i18n.T("error_fetching_data", lang))
				return
			}
			if err := b.DB.SetAlertChannel(i.GuildID, channelID); err != nil {
				b.respond(i, i18n.T("error_fetching_data", lang))
				return
			}
			b.respond(i, fmt.Sprintf("%s <#%s>", i18n.T("alert_channel_set", lang), channelID))
			return
		}
	}
	b.respond(i, i18n.T("error_fetching_data", lang))
}

func (b *Bot) handleSetupTideStation(i *discordgo.InteractionCreate, lang i18n.Language, op *discordgo.ApplicationCommandInteractionDataOption) {
	station := ""
	for _, opt := range op.Options {
		if opt.Name == "station" {
			station = opt.StringValue()
		}
	}
	station = strings.ToUpper(strings.TrimSpace(station))
	if !hko.IsValidTideStation(station) {
		b.respond(i, i18n.T("invalid_tide_station", lang))
		return
	}
	if err := b.DB.SetTideStation(i.GuildID, station); err != nil {
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}
	b.respond(i, fmt.Sprintf("%s **%s**.", i18n.T("tide_station_set", lang), station))
}

func (b *Bot) handleSetupStatus(i *discordgo.InteractionCreate, lang i18n.Language, op *discordgo.ApplicationCommandInteractionDataOption) {
	enabled := false
	for _, opt := range op.Options {
		if opt.Name == "enabled" {
			enabled = opt.BoolValue()
		}
	}
	if err := b.DB.SetBotStatusEnabled(i.GuildID, enabled); err != nil {
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}
	status := i18n.T("disabled", lang)
	if enabled {
		status = i18n.T("enabled", lang)
	}
	b.respond(i, fmt.Sprintf("%s **%s**.", i18n.T("status_set", lang), status))
}
