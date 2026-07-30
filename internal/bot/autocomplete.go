package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"weather-bot/internal/hko"
)

var warningCodes = []string{
	"WTCSGNL", "WRAIN", "WTS", "WMONSOON", "WFIRE",
	"WHOT", "WCOLD", "WFNT", "WHOTT", "WNTCC",
	"WRAINAMBER", "WRAINRED", "WRAINBLACK",
	"WTEMPERATURE",
}

func (b *Bot) handleAutocomplete(i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	for _, opt := range data.Options {
		if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			for _, subOpt := range opt.Options {
				if !subOpt.Focused {
					continue
				}
				switch subOpt.Name {
				case "type":
					b.respondWarningAutocomplete(i, subOpt.StringValue())
				case "station":
					b.respondStationAutocomplete(i, subOpt.StringValue())
				}
				return
			}
		}
	}
}

func (b *Bot) respondWarningAutocomplete(i *discordgo.InteractionCreate, input string) {
	input = strings.ToUpper(input)
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, code := range warningCodes {
		if strings.HasPrefix(code, input) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  code,
				Value: code,
			})
		}
		if len(choices) >= 25 {
			break
		}
	}
	_ = b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func (b *Bot) respondStationAutocomplete(i *discordgo.InteractionCreate, input string) {
	input = strings.ToUpper(input)
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, station := range hko.ValidTideStations() {
		if strings.HasPrefix(station, input) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  station,
				Value: station,
			})
		}
		if len(choices) >= 25 {
			break
		}
	}
	_ = b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}
