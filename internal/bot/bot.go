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
	b.Session.AddHandler(func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		defer func() {
			if r := recover(); r != nil {
				b.Logger.Error("handler panic recovered", slog.Any("panic", r))
				_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "An internal error occurred.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		}()
		b.handleInteraction(s, ic)
	})

	b.Session.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		if err := b.DB.DeleteGuildSettings(g.ID); err != nil {
			b.Logger.Warn("failed to clean up guild settings", slog.String("guild_id", g.ID), slog.Any("err", err))
		}
	})
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

// deferRespond sends a deferred response to buy time for slow commands.
func (b *Bot) deferRespond(i *discordgo.InteractionCreate) error {
	return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

// followUp edits the deferred response with embeds.
func (b *Bot) followUp(i *discordgo.InteractionCreate, embeds []*discordgo.MessageEmbed) {
	_, err := b.Session.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
	if err != nil {
		b.Logger.Error("failed to edit follow-up", slog.Any("err", err))
	}
}

// followUpError edits the deferred response with an error message.
func (b *Bot) followUpError(i *discordgo.InteractionCreate, content string) {
	_, err := b.Session.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		b.Logger.Error("failed to edit follow-up error", slog.Any("err", err))
	}
}

// hasManageServerPermission checks if the interaction user has Manage Server permission.
func (b *Bot) hasManageServerPermission(i *discordgo.InteractionCreate) bool {
	if i.Member == nil {
		return false
	}
	return i.Member.Permissions&discordgo.PermissionManageServer != 0
}

// respondComponentEphemeral sends an ephemeral error message for component interactions.
func (b *Bot) respondComponentEphemeral(i *discordgo.InteractionCreate, content string) {
	if err := b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		b.Logger.Error("failed to respond component ephemeral", slog.Any("err", err))
	}
}
