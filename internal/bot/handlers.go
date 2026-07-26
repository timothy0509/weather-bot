package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"

	"weather-bot/internal/hko"
	"weather-bot/internal/i18n"
)

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()

	switch data.Name {
	case "weather":
		b.handleWeatherCommand(i, data)
	case "setup":
		b.handleSetupCommand(i, data)
	}
}

func (b *Bot) handleWeatherCommand(i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	if len(data.Options) == 0 {
		b.respond(i, "Unknown command")
		return
	}
	op := data.Options[0]
	lang := b.guildLanguage(i.GuildID)

	switch op.Name {
	case "current":
		b.handleCurrent(i, lang)
	case "forecast":
		b.handleForecast(i, lang, 3)
	case "forecast9":
		b.handleForecast(i, lang, 9)
	case "warnings":
		b.handleWarnings(i, lang)
	case "warning":
		b.handleWarningDetail(i, lang, op)
	case "rain":
		b.handleRain(i, lang)
	case "uv":
		b.handleUV(i, lang)
	case "tide":
		b.handleTide(i, lang, op)
	case "lunar":
		b.handleLunar(i, lang)
	case "earthquake":
		b.handleEarthquake(i, lang)
	default:
		b.respond(i, "Unknown command")
	}
}

func (b *Bot) handleCurrent(i *discordgo.InteractionCreate, lang i18n.Language) {
	w, err := b.HKO.GetCurrentWeather(string(lang))
	if err != nil {
		b.Logger.Error("current weather failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("current_weather", lang),
		Color: 0x3498DB,
	}

	if hkoReading, ok := hko.ReadingByPlace(w.Temperature, "Hong Kong Observatory"); ok {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("temperature", lang),
			Value:  hkoReading.FormatTemperature(),
			Inline: true,
		})
	} else if len(w.Temperature.Data) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("temperature", lang),
			Value:  w.Temperature.Data[0].FormatTemperature(),
			Inline: true,
		})
	}

	if hkoReading, ok := hko.ReadingByPlace(w.Humidity, "Hong Kong Observatory"); ok {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("humidity", lang),
			Value:  fmt.Sprintf("%.0f%%", hkoReading.Value),
			Inline: true,
		})
	} else if len(w.Humidity.Data) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("humidity", lang),
			Value:  fmt.Sprintf("%.0f%%", w.Humidity.Data[0].Value),
			Inline: true,
		})
	}

	if max, ok := hko.MaxRainfallReading(w.Rainfall); ok {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("rain", lang),
			Value:  max.FormatRainfall(),
			Inline: true,
		})
	}

	if hkoReading, ok := hko.ReadingByPlace(w.UVIndex, "King's Park"); ok {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("uv_index", lang),
			Value:  fmt.Sprintf("%.1f", hkoReading.Value),
			Inline: true,
		})
	} else if len(w.UVIndex.Data) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   i18n.T("uv_index", lang),
			Value:  fmt.Sprintf("%.1f", w.UVIndex.Data[0].Value),
			Inline: true,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("%s: %s", i18n.T("updated_at", lang), w.UpdateTime),
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleForecast(i *discordgo.InteractionCreate, lang i18n.Language, days int) {
	f, err := b.HKO.GetForecast(string(lang))
	if err != nil {
		b.Logger.Error("forecast failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	limit := days
	if limit > len(f.WeatherForecast) {
		limit = len(f.WeatherForecast)
	}

	title := i18n.T("forecast", lang)
	if days == 9 {
		title = i18n.T("forecast_9day", lang)
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: f.GeneralSituation,
		Color:       0x2ECC71,
	}

	for _, day := range f.WeatherForecast[:limit] {
		maxTemp := ""
		minTemp := ""
		if day.ForecastMaxtemp.Unit == "C" {
			maxTemp = fmt.Sprintf("%.0f°C", day.ForecastMaxtemp.Value)
		} else {
			maxTemp = fmt.Sprintf("%.0f %s", day.ForecastMaxtemp.Value, day.ForecastMaxtemp.Unit)
		}
		if day.ForecastMintemp.Unit == "C" {
			minTemp = fmt.Sprintf("%.0f°C", day.ForecastMintemp.Value)
		} else {
			minTemp = fmt.Sprintf("%.0f %s", day.ForecastMintemp.Value, day.ForecastMintemp.Unit)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: day.Week,
			Value: fmt.Sprintf("%s: %s\n%s: %s\n%s: %s / %s\n%s: %s", 
				i18n.T("weather", lang), day.ForecastWeather,
				i18n.T("wind", lang), day.ForecastWind,
				i18n.T("temperature", lang), minTemp, maxTemp,
				i18n.T("psr", lang), day.PSR),
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("%s: %s", i18n.T("updated_at", lang), f.UpdateTime),
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleWarnings(i *discordgo.InteractionCreate, lang i18n.Language) {
	ws, err := b.HKO.GetWarningSummary(string(lang))
	if err != nil {
		b.Logger.Error("warning summary failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("warnings", lang),
		Color: 0xE74C3C,
	}

	if len(ws) == 0 {
		embed.Description = i18n.T("no_active_warnings", lang)
	} else {
		for code, w := range ws {
			action := w.ActionCode
			if action == "" {
				action = "Active"
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s (%s)", w.Name, code),
				Value:  fmt.Sprintf("%s: %s\n%s", w.Type, action, w.UpdateTime),
				Inline: true,
			})
		}
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleWarningDetail(i *discordgo.InteractionCreate, lang i18n.Language, op *discordgo.ApplicationCommandInteractionDataOption) {
	warningType := ""
	for _, opt := range op.Options {
		if opt.Name == "type" {
			warningType = opt.StringValue()
		}
	}
	if warningType == "" {
		b.respond(i, "Please provide a warning type.")
		return
	}

	info, err := b.HKO.GetWarningInfo(string(lang))
	if err != nil {
		b.Logger.Error("warning info failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("warning_detail", lang),
		Color: 0xE74C3C,
	}

	found := false
	for _, d := range info.Details {
		if d.WarningStatementCode != warningType {
			continue
		}
		found = true
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   d.WarningStatementCode,
			Value:  strings.Join(d.Contents, "\n"),
			Inline: false,
		})
	}
	if !found {
		embed.Description = i18n.T("no_active_warnings", lang)
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleRain(i *discordgo.InteractionCreate, lang i18n.Language) {
	r, err := b.HKO.GetHourlyRainfall(string(lang))
	if err != nil {
		b.Logger.Error("rainfall failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("rainfall", lang),
		Color: 0x3498DB,
	}

	if len(r.HourlyRainfall) == 0 || len(r.HourlyRainfall[0].Data) == 0 {
		embed.Description = i18n.T("no_data", lang)
	} else {
		var sb strings.Builder
		group := r.HourlyRainfall[0]
		for _, reading := range group.Data {
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("**%s**: %.1f %s", reading.Place, reading.Max, reading.Unit))
		}
		embed.Description = sb.String()
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleUV(i *discordgo.InteractionCreate, lang i18n.Language) {
	w, err := b.HKO.GetCurrentWeather(string(lang))
	if err != nil {
		b.Logger.Error("uv failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("uv_index", lang),
		Color: 0xF1C40F,
	}

	if r, ok := hko.ReadingByPlace(w.UVIndex, "King's Park"); ok {
		embed.Description = fmt.Sprintf("%s: %.1f", r.Place, r.Value)
	} else if len(w.UVIndex.Data) > 0 {
		embed.Description = fmt.Sprintf("%s: %.1f", w.UVIndex.Data[0].Place, w.UVIndex.Data[0].Value)
	} else {
		embed.Description = i18n.T("no_data", lang)
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleTide(i *discordgo.InteractionCreate, lang i18n.Language, op *discordgo.ApplicationCommandInteractionDataOption) {
	station := b.guildTideStation(i.GuildID)
	for _, opt := range op.Options {
		if opt.Name == "station" {
			station = strings.ToUpper(opt.StringValue())
		}
	}
	if !hko.IsValidTideStation(station) {
		b.respond(i, i18n.T("invalid_tide_station", lang))
		return
	}

	tides, err := b.HKO.GetTides(station, string(lang))
	if err != nil {
		b.Logger.Error("tides failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (%s)", i18n.T("tides", lang), station),
		Color: 0x1ABC9C,
	}

	var sb strings.Builder
	for _, rec := range tides.Records() {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%02d:00: %.2f m", rec.Hour, rec.Height))
	}
	if sb.Len() == 0 {
		embed.Description = i18n.T("no_data", lang)
	} else {
		embed.Description = sb.String()
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleLunar(i *discordgo.InteractionCreate, lang i18n.Language) {
	l, err := b.HKO.GetTodayLunarCalendar()
	if err != nil {
		b.Logger.Error("lunar failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("lunar_calendar", lang),
		Color: 0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Year / 年",
				Value:  l.LunarYear,
				Inline: true,
			},
			{
				Name:   "Date / 日期",
				Value:  l.LunarDate,
				Inline: true,
			},
		},
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func (b *Bot) handleEarthquake(i *discordgo.InteractionCreate, lang i18n.Language) {
	e, err := b.HKO.GetEarthquakeInfo(string(lang))
	if err != nil {
		b.Logger.Error("earthquake failed", slog.Any("err", err))
		b.respond(i, i18n.T("error_fetching_data", lang))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: i18n.T("earthquake", lang),
		Color: 0xE67E22,
	}

	results := e.Results()
	if len(results) == 0 {
		embed.Description = i18n.T("no_earthquake", lang)
	} else {
		for _, q := range results[:min(3, len(results))] {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s %s", i18n.T("magnitude", lang), q.FormatMag()),
				Value:  fmt.Sprintf("%s: %s\n%s: %s", i18n.T("region", lang), q.Region, i18n.T("time", lang), q.PTime),
				Inline: true,
			})
		}
	}
	b.respondEmbed(i, []*discordgo.MessageEmbed{embed})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
