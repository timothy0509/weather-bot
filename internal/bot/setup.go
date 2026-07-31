package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"

	"weather-bot/internal/hko"
	"weather-bot/internal/i18n"
)

func (b *Bot) handleSetupCommand(i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	lang := b.guildLanguage(i.GuildID)

	if !b.hasManageServerPermission(i) {
		b.respondComponentEphemeral(i, i18n.T("setup_permission_required", lang))
		return
	}

	b.respondSetupMainPanel(i, lang)
}

func (b *Bot) handleSetupComponent(i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	customID := data.CustomID

	if !strings.HasPrefix(customID, "setup:") {
		return
	}

	if !b.hasManageServerPermission(i) {
		b.respondComponentEphemeral(i, i18n.T("setup_permission_required", b.guildLanguage(i.GuildID)))
		return
	}

	lang := b.guildLanguage(i.GuildID)

	switch customID {
	case "setup:main":
		b.editSetupMainPanel(i, lang)
	case "setup:lang":
		b.editSetupLanguagePanel(i, lang)
	case "setup:lang:set":
		if len(data.Values) > 0 {
			b.handleSetupLangSet(i, lang, data.Values[0])
		}
	case "setup:alert":
		b.editSetupAlertPanel(i, lang)
	case "setup:alert:set":
		b.handleSetupAlertSet(i, lang, data)
	case "setup:alert:disable":
		b.handleSetupAlertDisable(i, lang)
	case "setup:tide":
		b.editSetupTidePanel(i, lang)
	case "setup:tide:set":
		if len(data.Values) > 0 {
			b.handleSetupTideSet(i, lang, data.Values[0])
		}
	case "setup:status":
		b.editSetupStatusPanel(i, lang)
	case "setup:status:set:true":
		b.handleSetupStatusSet(i, lang, true)
	case "setup:status:set:false":
		b.handleSetupStatusSet(i, lang, false)
	}
}

func (b *Bot) respondSetupMainPanel(i *discordgo.InteractionCreate, lang i18n.Language) {
	embed, components := b.buildMainPanel(i.GuildID, lang)
	if err := b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		b.Logger.Error("failed to respond setup panel", slog.Any("err", err))
	}
}

func (b *Bot) editSetupMainPanel(i *discordgo.InteractionCreate, lang i18n.Language) {
	embed, components := b.buildMainPanel(i.GuildID, lang)
	b.editSetupPanel(i, embed, components)
}

func (b *Bot) editSetupLanguagePanel(i *discordgo.InteractionCreate, lang i18n.Language) {
	embed, components := b.buildLanguagePanel(i.GuildID, lang)
	b.editSetupPanel(i, embed, components)
}

func (b *Bot) editSetupAlertPanel(i *discordgo.InteractionCreate, lang i18n.Language) {
	embed, components := b.buildAlertPanel(i.GuildID, lang)
	b.editSetupPanel(i, embed, components)
}

func (b *Bot) editSetupTidePanel(i *discordgo.InteractionCreate, lang i18n.Language) {
	embed, components := b.buildTidePanel(i.GuildID, lang)
	b.editSetupPanel(i, embed, components)
}

func (b *Bot) editSetupStatusPanel(i *discordgo.InteractionCreate, lang i18n.Language) {
	embed, components := b.buildStatusPanel(i.GuildID, lang)
	b.editSetupPanel(i, embed, components)
}

func (b *Bot) editSetupPanel(i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	if err := b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	}); err != nil {
		b.Logger.Error("failed to edit setup panel", slog.Any("err", err))
	}
}

func (b *Bot) buildMainPanel(guildID string, lang i18n.Language) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		b.Logger.Error("failed to get guild settings", slog.Any("err", err))
		return nil, nil
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

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    i18n.T("btn_language", lang),
					CustomID: "setup:lang",
					Style:    discordgo.SecondaryButton,
				},
				discordgo.Button{
					Label:    i18n.T("btn_alert_channel", lang),
					CustomID: "setup:alert",
					Style:    discordgo.SecondaryButton,
				},
				discordgo.Button{
					Label:    i18n.T("btn_tide_station", lang),
					CustomID: "setup:tide",
					Style:    discordgo.SecondaryButton,
				},
				discordgo.Button{
					Label:    i18n.T("btn_bot_status", lang),
					CustomID: "setup:status",
					Style:    discordgo.SecondaryButton,
				},
			},
		},
	}

	return embed, components
}

func (b *Bot) buildLanguagePanel(guildID string, lang i18n.Language) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		b.Logger.Error("failed to get guild settings", slog.Any("err", err))
		return nil, nil
	}

	embed := &discordgo.MessageEmbed{
		Title:       i18n.T("btn_language", lang),
		Description: fmt.Sprintf("**%s:** %s", i18n.T("current_setting", lang), gs.Language),
		Color:       0x3498DB,
	}

	options := []discordgo.SelectMenuOption{
		{
			Label:   "English",
			Value:   "en",
			Default: gs.Language == "en",
		},
		{
			Label:   "繁體中文",
			Value:   "tc",
			Default: gs.Language == "tc",
		},
		{
			Label:   "简体中文",
			Value:   "sc",
			Default: gs.Language == "sc",
		},
		{
			Label:   "Bilingual",
			Value:   "bilingual",
			Default: gs.Language == "bilingual",
		},
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:    discordgo.StringSelectMenu,
					CustomID:    "setup:lang:set",
					Placeholder: i18n.T("select_language_prompt", lang),
					Options:     options,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    i18n.T("btn_back", lang),
					CustomID: "setup:main",
					Style:    discordgo.SecondaryButton,
				},
			},
		},
	}

	return embed, components
}

func (b *Bot) buildAlertPanel(guildID string, lang i18n.Language) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		b.Logger.Error("failed to get guild settings", slog.Any("err", err))
		return nil, nil
	}

	channelValue := i18n.T("alert_channel_disabled", lang)
	if gs.AlertChannelID.Valid {
		channelValue = fmt.Sprintf("<#%s>", gs.AlertChannelID.String)
	}

	embed := &discordgo.MessageEmbed{
		Title:       i18n.T("btn_alert_channel", lang),
		Description: fmt.Sprintf("**%s:** %s", i18n.T("current_setting", lang), channelValue),
		Color:       0xE67E22,
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:     discordgo.ChannelSelectMenu,
					CustomID:     "setup:alert:set",
					Placeholder:  i18n.T("select_channel_prompt", lang),
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    i18n.T("btn_disable", lang),
					CustomID: "setup:alert:disable",
					Style:    discordgo.DangerButton,
				},
				discordgo.Button{
					Label:    i18n.T("btn_back", lang),
					CustomID: "setup:main",
					Style:    discordgo.SecondaryButton,
				},
			},
		},
	}

	return embed, components
}

func (b *Bot) buildTidePanel(guildID string, lang i18n.Language) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		b.Logger.Error("failed to get guild settings", slog.Any("err", err))
		return nil, nil
	}

	embed := &discordgo.MessageEmbed{
		Title:       i18n.T("btn_tide_station", lang),
		Description: fmt.Sprintf("**%s:** %s", i18n.T("current_setting", lang), gs.TideStation),
		Color:       0x1ABC9C,
	}

	var options []discordgo.SelectMenuOption
	for _, station := range hko.ValidTideStations() {
		options = append(options, discordgo.SelectMenuOption{
			Label:   station,
			Value:   station,
			Default: gs.TideStation == station,
		})
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:    discordgo.StringSelectMenu,
					CustomID:    "setup:tide:set",
					Placeholder: i18n.T("select_tide_prompt", lang),
					Options:     options,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    i18n.T("btn_back", lang),
					CustomID: "setup:main",
					Style:    discordgo.SecondaryButton,
				},
			},
		},
	}

	return embed, components
}

func (b *Bot) buildStatusPanel(guildID string, lang i18n.Language) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		b.Logger.Error("failed to get guild settings", slog.Any("err", err))
		return nil, nil
	}

	statusValue := i18n.T("disabled", lang)
	if gs.BotStatusEnabled {
		statusValue = i18n.T("enabled", lang)
	}

	embed := &discordgo.MessageEmbed{
		Title:       i18n.T("btn_bot_status", lang),
		Description: fmt.Sprintf("**%s:** %s", i18n.T("current_setting", lang), statusValue),
		Color:       0x9B59B6,
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    i18n.T("btn_enable", lang),
					CustomID: "setup:status:set:true",
					Style:    discordgo.SuccessButton,
				},
				discordgo.Button{
					Label:    i18n.T("btn_disable", lang),
					CustomID: "setup:status:set:false",
					Style:    discordgo.DangerButton,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    i18n.T("btn_back", lang),
					CustomID: "setup:main",
					Style:    discordgo.SecondaryButton,
				},
			},
		},
	}

	return embed, components
}

func (b *Bot) handleSetupLangSet(i *discordgo.InteractionCreate, lang i18n.Language, mode string) {
	if !i18n.IsValid(mode) {
		b.respondComponentEphemeral(i, i18n.T("invalid_language_mode", lang))
		return
	}
	if err := b.DB.SetLanguage(i.GuildID, mode); err != nil {
		b.respondComponentEphemeral(i, i18n.T("error_fetching_data", lang))
		return
	}
	newLang := i18n.Normalize(mode)
	b.editSetupMainPanel(i, newLang)
}

func (b *Bot) handleSetupAlertSet(i *discordgo.InteractionCreate, lang i18n.Language, data discordgo.MessageComponentInteractionData) {
	if len(data.Resolved.Channels) == 0 {
		b.respondComponentEphemeral(i, i18n.T("error_fetching_data", lang))
		return
	}

	var channelID string
	for id := range data.Resolved.Channels {
		channelID = id
		break
	}

	if err := b.DB.SetAlertChannel(i.GuildID, channelID); err != nil {
		b.respondComponentEphemeral(i, i18n.T("error_fetching_data", lang))
		return
	}
	b.editSetupMainPanel(i, lang)
}

func (b *Bot) handleSetupAlertDisable(i *discordgo.InteractionCreate, lang i18n.Language) {
	if err := b.DB.SetAlertChannel(i.GuildID, ""); err != nil {
		b.respondComponentEphemeral(i, i18n.T("error_fetching_data", lang))
		return
	}
	b.editSetupMainPanel(i, lang)
}

func (b *Bot) handleSetupTideSet(i *discordgo.InteractionCreate, lang i18n.Language, station string) {
	station = strings.ToUpper(strings.TrimSpace(station))
	if !hko.IsValidTideStation(station) {
		b.respondComponentEphemeral(i, i18n.T("invalid_tide_station", lang))
		return
	}
	if err := b.DB.SetTideStation(i.GuildID, station); err != nil {
		b.respondComponentEphemeral(i, i18n.T("error_fetching_data", lang))
		return
	}
	b.editSetupMainPanel(i, lang)
}

func (b *Bot) handleSetupStatusSet(i *discordgo.InteractionCreate, lang i18n.Language, enabled bool) {
	if err := b.DB.SetBotStatusEnabled(i.GuildID, enabled); err != nil {
		b.respondComponentEphemeral(i, i18n.T("error_fetching_data", lang))
		return
	}
	b.editSetupMainPanel(i, lang)
}
