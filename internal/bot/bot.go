// Package bot implements Discord bot logic for the Hong Kong weather bot.
package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"weather-bot/internal/config"
	"weather-bot/internal/db"
	"weather-bot/internal/hko"
	"weather-bot/internal/i18n"
)

// Bot is the Discord bot instance.
type Bot struct {
	Session *discordgo.Session
	Config  *config.Config
	DB      *db.DB
	HKO     *hko.Client
	Logger  *slog.Logger
}

// New creates a new Bot instance.
func New(cfg *config.Config, database *db.DB, logger *slog.Logger) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	return &Bot{
		Session: session,
		Config:  cfg,
		DB:      database,
		HKO:     hko.NewClient(0),
		Logger:  logger,
	}, nil
}

// Close closes the bot session.
func (b *Bot) Close() error {
	return b.Session.Close()
}

// RegisterCommands registers slash commands for the bot.
func (b *Bot) RegisterCommands() error {
	return registerCommands(b.Session, b.Config.ApplicationID, b.Config.GuildID)
}

// RegisterHandlers registers Discord interaction handlers.
func (b *Bot) RegisterHandlers() {
	b.Session.AddHandler(b.handleInteraction)
}

// guildLanguage returns the configured language for a guild.
func (b *Bot) guildLanguage(guildID string) i18n.Language {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		b.Logger.Warn("failed to get guild settings", slog.String("guild_id", guildID), slog.Any("err", err))
		return i18n.EN
	}
	return i18n.Normalize(gs.Language)
}

// guildTideStation returns the configured tide station for a guild.
func (b *Bot) guildTideStation(guildID string) string {
	gs, err := b.DB.GetGuildSettings(guildID)
	if err != nil {
		return "QUB"
	}
	return gs.TideStation
}

// respond sends a simple interaction response.
func (b *Bot) respond(i *discordgo.InteractionCreate, content string) {
	if err := b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}); err != nil {
		b.Logger.Error("failed to respond", slog.Any("err", err))
	}
}

// respondEmbed sends an embed interaction response.
func (b *Bot) respondEmbed(i *discordgo.InteractionCreate, embeds []*discordgo.MessageEmbed) {
	if err := b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	}); err != nil {
		b.Logger.Error("failed to respond embed", slog.Any("err", err))
	}
}
