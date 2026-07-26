# Hong Kong Observatory Open Data API Reference

> For use with a Discord bot collecting live weather information.
> Base documentation: `HKO_Open_Data_API_Documentation.pdf` v1.13 (Sep 2025)
> No API key required. All endpoints are free and public.

---

## Table of Contents

1. [Weather Information API](#1-weather-information-api) — forecasts, current weather, warnings
2. [Earthquake Information API](#2-earthquake-information-api) — seismic events
3. [Open Data (Climate & Weather) API](#3-open-data-climate-and-weather-information-api) — tides, lightning, visibility, historical climate
4. [Gregorian-Lunar Calendar Conversion API](#4-gregorian-lunar-calendar-conversion-api)
5. [Rainfall in the Past Hour API](#5-rainfall-in-the-past-hour-from-automatic-weather-station-api)
6. [RSS Feeds and Other Live Data Sources](#6-rss-feeds-and-other-live-data-sources) — RSS, XML, imagery, radar, satellite
7. [Discord Bot Implementation Patterns](#7-discord-bot-implementation-patterns)

---

## 1. Weather Information API

**Endpoint:** `https://data.weather.gov.hk/weatherAPI/opendata/weather.php`
**Method:** GET
**Return Type:** JSON
**Language:** `lang=en|tc|sc` (default: en)

### 1.1 `dataType=flw` — Local Weather Forecast

Plain-text forecast with tropical cyclone information.

**Example:** `?dataType=flw&lang=en`

```json
{
    "generalSituation": "Tropical Cyclone Noul gradually moved towards...",
    "tcInfo": "",
    "fireDangerWarning": "",
    "forecastPeriod": "Weather forecast for Hong Kong",
    "forecastDesc": "Gale force northwesterly winds...",
    "outlook": "There will still be showers on Monday...",
    "updateTime": "2026-07-26T00:45:00+08:00"
}
```

| Field | Description |
|---|---|
| `generalSituation` | General weather situation text |
| `tcInfo` | Tropical cyclone information (empty if none active) |
| `fireDangerWarning` | Fire danger warning message |
| `forecastPeriod` | "Weather forecast for Hong Kong" |
| `forecastDesc` | Detailed forecast description |
| `outlook` | Extended outlook text |
| `updateTime` | ISO 8601 timestamp (GMT+8) |

### 1.2 `dataType=fnd` — 9-day Weather Forecast

Structured daily forecast with numerical data.

**Example:** `?dataType=fnd&lang=en`

```json
{
    "generalSituation": "...",
    "weatherForecast": [
        {
            "forecastDate": "20260726",
            "week": "Sunday",
            "forecastWind": "Northwest force 8 to 9...",
            "forecastWeather": "Cloudy with frequent heavy squally showers...",
            "forecastMaxtemp": { "value": 29, "unit": "C" },
            "forecastMintemp": { "value": 26, "unit": "C" },
            "forecastMaxrh": { "value": 95, "unit": "percent" },
            "forecastMinrh": { "value": 80, "unit": "percent" },
            "ForecastIcon": 64,
            "PSR": "High"
        }
    ],
    "updateTime": "2026-07-26T00:45:00+08:00",
    "seaTemp": [ ... ],
    "soilTemp": [ ... ]
}
```

| Field | Description |
|---|---|
| `weatherForecast[]` | Array of 9 daily forecasts |
| `forecastDate` | Date in `YYYYMMDD` |
| `week` | Day of week |
| `forecastWind` | Wind description |
| `forecastWeather` | Weather description |
| `forecastMaxtemp` / `forecastMintemp` | `{ value, unit }` |
| `forecastMaxrh` / `forecastMinrh` | Relative humidity `{ value, unit }` |
| `ForecastIcon` | Weather icon number (see [icon list](https://www.hko.gov.hk/textonly/v2/explain/wxicon_e.htm)) |
| `PSR` | Probability of Significant Rain: `High` / `Medium High` / `Medium` / `Medium Low` / `Low` |
| `seaTemp` | Array of sea temperature readings at stations |
| `soilTemp` | Array of soil temperature readings at stations |

**Discord bot use:** Display the first 1-3 days in an embed, or show all 9 in a paginated view.

### 1.3 `dataType=rhrread` — Current Weather Report

Real-time observations from stations across Hong Kong.

**Example:** `?dataType=rhrread&lang=en`

Top-level keys returned:
- `temperature` — Array of `{ place, value, unit }` readings
- `humidity` — `{ recordTime, data: [{ place, value, unit }] }`
- `rainfall` — District rainfall data with `{ place, max, min, unit, main }`
- `icon` — Weather icon codes (array)
- `iconUpdateTime` — Timestamp
- `uvindex` — UV index data with `{ place, value, desc }`
- `specialWxTips` — Special weather tip messages
- `warningMessage` — Active warning messages
- `tcmessage` — Tropical cyclone position message
- `mintempFrom00To09` — Minimum temp from midnight to 9am
- `rainfallFrom00To12` — Accumulated rainfall midnight to noon
- `rainfallLastMonth` — Rainfall in previous month
- `rainfallJanuaryToLastMonth` — Accumulated rainfall since January
- `updateTime` — Timestamp

**Discord bot use:** Primary endpoint for current conditions. Show temperature at HKO or Kings Park, humidity, UV index, and rainfall.

### 1.4 `dataType=warnsum` — Weather Warning Summary

Current active warning signals with status codes.

**Example:** `?dataType=warnsum&lang=en`

```json
{
    "WTCSGNL": {
        "name": "Tropical Cyclone Warning Signal",
        "code": "TC8NW",
        "type": "No. 8 Northwest Gale or Storm Signal",
        "actionCode": "ISSUE",
        "issueTime": "2026-07-26T00:15:00+08:00",
        "updateTime": "2026-07-26T00:15:00+08:00"
    }
}
```

Warning codes (keys):
| Code | Warning |
|---|---|
| `WFIRE` | Fire Danger Warning |
| `WFROST` | Frost Warning |
| `WHOT` | Hot Weather Warning |
| `WCOLD` | Cold Weather Warning |
| `WMSGNL` | Strong Monsoon Signal |
| `WRAIN` | Rainstorm Warning Signal |
| `WFNTSA` | Special Announcement on Flooding in northern NT |
| `WL` | Landslip Warning |
| `WTCSGNL` | Tropical Cyclone Warning Signal |
| `WTMW` | Tsunami Warning |
| `WTS` | Thunderstorm Warning |

Signal subtypes (within `code` field):
- TC: `TC1`, `TC3`, `TC8NE`, `TC8SE`, `TC8SW`, `TC8NW`, `TC9`, `TC10`, `CANCEL`
- Rainstorm: `WRAINA` (Amber), `WRAINR` (Red), `WRAINB` (Black)
- Fire: `WFIREY` (Yellow), `WFIRER` (Red)

`actionCode`: `ISSUE`, `REISSUE`, `CANCEL`, `EXTEND`, `UPDATE`

**Discord bot use:** Alert channel for when typhoon signals, rainstorms, or thunderstorm warnings are active.

### 1.5 `dataType=warningInfo` — Weather Warning Information

Detailed bulletins for each active warning, with full text content.

**Example:** `?dataType=warningInfo&lang=en`

```json
{
    "details": [
        {
            "contents": ["Warning text line 1...", "Warning text line 2..."],
            "warningStatementCode": "WTCSGNL",
            "subtype": "TC8NW",
            "updateTime": "2026-07-26T00:15:00+08:00"
        }
    ]
}
```

`warningStatementCode` uses same codes as `warnsum`.
`subtype` only present for fire danger, tropical cyclone, and rainstorm warnings.

**Discord bot use:** When a warning is active, fetch this for the detailed text to display.

### 1.6 `dataType=swt` — Special Weather Tips

Ad-hoc weather tips (e.g., pre-no.8 announcements, localized heavy rain).

```json
{
    "swt": [
        {
            "desc": "The Hong Kong Observatory announces that...",
            "updateTime": "2026-07-26T14:10:00+08:00"
        }
    ]
}
```

**Discord bot use:** Check for non-null/non-empty to display important announcements.

---

## 2. Earthquake Information API

**Endpoint:** `https://data.weather.gov.hk/weatherAPI/opendata/earthquake.php`
**Method:** GET
**Return Type:** JSON

### 2.1 `dataType=qem` — Quick Earthquake Messages

Recent earthquakes detected globally.

**Example:** `?dataType=qem&lang=en`

| Field | Description |
|---|---|
| `lat` | Latitude |
| `lon` | Longitude |
| `mag` | Richter magnitude |
| `region` | Region description |
| `ptime` | Earthquake date/time (ISO 8601) |
| `updateTime` | Last update time |

### 2.2 `dataType=feltearthquake` — Locally Felt Earth Tremor Report

Earthquakes felt in Hong Kong.

| Field | Description |
|---|---|
| `mag` | Richter magnitude |
| `region` | Region |
| `intensity` | Intensity felt in HK |
| `lat` / `lon` | Coordinates |
| `details` | Description |
| `ptime` | Earthquake time |
| `updateTime` | Last update |

**Discord bot use:** Low priority for a weather bot, but could be used for a natural-disaster alert channel.

---

## 3. Open Data (Climate and Weather Information) API

**Endpoint:** `https://data.weather.gov.hk/weatherAPI/opendata/opendata.php`
**Method:** GET
**Return Type:** JSON or CSV (`rformat=json|csv`, default CSV)

### 3.1 `dataType=HHOT` — Hourly Astronomical Tide Heights

**Parameters:** `station`, `year` (2022-2024), `month` (1-12 opt), `day` (1-31 opt), `hour` (1-24 opt)

Station codes: `CCH` (Cheung Chau), `CLK` (Chek Lap Kok), `CMW` (Chi Ma Wan), `KCT` (Kwai Chung), `KLW` (Ko Lau Wan), `LOP` (Lok On Pai), `MWC` (Ma Wan), `QUB` (Quarry Bay), `SPW` (Shek Pik), `TAO` (Tai O), `TBT` (Tsim Bei Tsui), `TMW` (Tai Miu Wan), `TPK` (Tai Po Kau), `WAG` (Waglan Island)

### 3.2 `dataType=HLT` — Times and Heights of High/Low Tides

Same station and date parameters as HHOT.

### 3.3 `dataType=SRS` — Times of Sunrise, Sun Transit and Sunset

**Parameters:** `year`, `month` (opt), `day` (opt)

### 3.4 `dataType=MRS` — Times of Moonrise, Moon Transit and Moonset

**Parameters:** `year` (2018-2024), `month` (opt), `day` (opt)

### 3.5 `dataType=LHL` — Cloud-to-ground and Cloud-to-cloud Lightning Count

Real-time lightning data.

**Parameters:** `lang=en|tc|sc`

### 3.6 `dataType=LTMV` — Latest 10-minute Mean Visibility

**Parameters:** `lang=en|tc|sc`

### 3.7 `dataType=CLMTEMP` — Daily Mean Temperature

**Parameters:** `station`, `year` (1884-current, varies by station), `month` (opt)

Station codes (partial list): `CCH`, `CWB`, `HKA`, `HKO`, `HKP`, `HKS`, `HPV`, `JKB`, `KLT`, `KP`, `KSC`, `KTG`, `LFS`, `NGP`, `PEN`, `PLC`, `SE1`, `SEK`, `SHA`, `SKG`, `SKW`, `SSH`, `SSP`, `STY`, `TC`, `TKL`, `TMS`, `TPO`, `TU1`, `TW`, `TWN`, `TY1`, `TYW`, `VP1`, `WGL`, `WLP`, `WTS`, `YCT`, `YLP`

### 3.8 `dataType=CLMMAXT` — Daily Maximum Temperature

Same parameters as CLMTEMP.

### 3.9 `dataType=CLMMINT` — Daily Minimum Temperature

Same parameters as CLMTEMP.

### 3.10 `dataType=RYES` — Weather and Radiation Level Report

**Parameters:** `date` (YYYYMMDD, 20190910 to yesterday), `lang=en|tc|sc`

---

## 4. Gregorian-Lunar Calendar Conversion API

**Endpoint:** `https://data.weather.gov.hk/weatherAPI/opendata/lunardate.php`
**Method:** GET
**Return Type:** JSON

**Parameters:** `date` in `YYYY-MM-DD` format (current year + 2 years)

**Example:** `?date=2026-07-26`

```json
{
    "LunarYear": "丙午年，馬",
    "LunarDate": "六月十二"
}
```

| Field | Description |
|---|---|
| `LunarYear` | Gan-Zhi year + zodiac (traditional Chinese only) |
| `LunarDate` | Lunar date (traditional Chinese only) |

---

## 5. RSS Feeds and Other Live Data Sources

### 5.1 RSS Feeds

**Base domain:** `https://rss.weather.gov.hk/rss/`
**Method:** GET
**Return Type:** RSS/XML (`application/rss+xml`)
**Languages:** English (`.xml`), Traditional Chinese (`_uc.xml`), Simplified Chinese (`sc/rss/`)

All RSS feeds follow standard RSS 2.0 with `<channel>`, `<title>`, `<link>`, `<description>`, `<item>`, `<pubDate>`, etc.

#### 5.1.1 Current Weather Report

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/CurrentWeather.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/CurrentWeather_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/CurrentWeather_uc.xml` |

**Update frequency:** Hourly and on significant change.

Each `<item>` contains a weather station report with temperature, humidity, wind, rainfall, and weather condition text.

**Discord bot use:** Poll as a lightweight alternative to the `rhrread` API. Parse `<item>` elements for individual station data.

#### 5.1.2 Local Weather Forecast

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/LocalWeatherForecast.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/LocalWeatherForecast_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/LocalWeatherForecast_uc.xml` |

**Update frequency:** Hourly and on significant change.

Plain-text forecast for today and tomorrow, plus tropical cyclone info if active. Lightweight alternative to `dataType=flw`.

#### 5.1.3 9-Day Weather Forecast

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/SeveralDaysWeatherForecast_v2.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/SeveralDaysWeatherForecast_v2_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/SeveralDaysWeatherForecast_v2_uc.xml` |

**Update frequency:** Twice daily and on significant change.

Each day is a separate `<item>` with forecast text. Lightweight alternative to `dataType=fnd`.

#### 5.1.4 Weather Warning Information

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/WeatherWarningBulletin.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/WeatherWarningBulletin_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/WeatherWarningBulletin_uc.xml` |

**Update frequency:** As needed when warnings change.

Full warning bulletin text. Similar content to `dataType=warningInfo`.

#### 5.1.5 Weather Warning Summary

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/WeatherWarningSummaryv2.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/WeatherWarningSummaryv2_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/WeatherWarningSummaryv2_uc.xml` |

**Update frequency:** As needed when warnings change.

Concise summary of all active warnings. Similar content to `dataType=warnsum`.

#### 5.1.6 Quick Earthquake Messages

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/QuickEarthquakeMessage.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/QuickEarthquakeMessage_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/QuickEarthquakeMessage_uc.xml` |

**Update frequency:** As needed when earthquakes M6.0+ occur.

#### 5.1.7 Locally Felt Earth Tremor Report

| Language | URL |
|---|---|
| English | `https://rss.weather.gov.hk/rss/FeltEarthquake.xml` |
| Trad. Chinese | `https://rss.weather.gov.hk/rss/FeltEarthquake_uc.xml` |
| Simp. Chinese | `https://rss.weather.gov.hk/sc/rss/FeltEarthquake_uc.xml` |

**Update frequency:** As needed when tremors felt in Hong Kong.

### 5.2 Live XML Feeds

#### 5.2.1 Tropical Cyclone Track Information

**URL:** `https://www.weather.gov.hk/wxinfo/currwx/tc_list.xml`
**Format:** XML
**Update frequency:** As needed during TC activity.

Contains real-time tropical cyclone track data: position, intensity, movement direction/speed, and forecast tracks.

### 5.3 Live Imagery Feeds

#### 5.3.1 Weather Radar Images

**Base URL:** `https://www.hko.gov.hk/wxinfo/radars/`
**Update interval:** Every 6 minutes
**Format:** Dynamic JPG images named by timestamp (`YYYYMMDDHHMM`)

| Range / Height | Path | Example |
|---|---|---|
| 256 km, 3 km | `rad_256_png/2d256nradar_{ts}.jpg` | `https://www.hko.gov.hk/wxinfo/radars/rad_256_png/2d256nradar_202607260100.jpg` |
| 128 km, 3 km | `rad_128_png/2d128nradar_{ts}.jpg` | `https://www.hko.gov.hk/wxinfo/radars/rad_128_png/2d128nradar_202607260100.jpg` |
| 64 km, 3 km | `rad_064_png/2d064nradar_{ts}.jpg` | `https://www.hko.gov.hk/wxinfo/radars/rad_064_png/2d064nradar_202607260100.jpg` |
| 64 km, 2 km | `rad_2km_064_png/2d064_2km_nradar_{ts}.jpg` | `https://www.hko.gov.hk/wxinfo/radars/rad_2km_064_png/2d064_2km_nradar_202607260100.jpg` |

**JSON index** (for discovering latest image filenames):
`https://www.hko.gov.hk/wxinfo/radars/temp_json/nradar_img.json`

**Interactive viewers:**
- 64 km: `https://www.hko.gov.hk/en/wxinfo/radars/radar_range1.htm`
- 128 km: `https://www.hko.gov.hk/en/wxinfo/radars/radar_range0.htm`
- 256 km: `https://www.hko.gov.hk/en/wxinfo/radars/radar.htm`
- Radar + lightning: `https://www.hko.gov.hk/en/wxinfo/llis/llisradar.shtml`

**Discord bot use:** Fetch latest radar image and post to a weather channel. Use the JSON index to discover current image filenames, or compute timestamps rounded to nearest 6 minutes.

#### 5.3.2 Weather Satellite Imagery

**Interactive viewer:** `https://www.hko.gov.hk/en/wxinfo/intersat/satellite/sate.htm`
**Image gallery:** `https://www.hko.gov.hk/en/wxinfo/intersat/satellite_gallery/index.htm`

Sources: Himawari-9, FY-2G
Types: Infra-red, True Colour, Deep Convection, Visible, Water Vapour

#### 5.3.3 Lightning Location Data

**Interactive map:** `https://www.hko.gov.hk/en/wxinfo/llis/gm_index.htm`
**Radar overlay:** `https://www.hko.gov.hk/en/wxinfo/llis/llisradar.shtml`
**Update interval:** Every 5 minutes

Shows cloud-to-ground and cloud-to-cloud lightning strokes in the past 30 minutes.

### 5.4 Other Live Data Pages (Web/Interactive)

| Data Type | URL | Description |
|---|---|---|
| Regional Weather Portal | `https://www.hko.gov.hk/en/wxinfo/awsgis/regional_portal.html` | Interactive map with temp, humidity, wind, pressure |
| Rainfall Distribution Map | `https://www.hko.gov.hk/en/wxinfo/rainfall/isohyet_past1hr.shtml` | Surface rainfall distribution over last hour |
| Automatic Regional Forecast | `https://maps.weather.gov.hk/ocf/index_e.html` | Gridded forecast for HK & Pearl River Delta |
| Coastal Sea Level (Tide) | `https://www.hko.gov.hk/en/tide/marine/realtide.htm` | Real-time tide/sea level data |
| Astronomical Tide Prediction | `https://www.hko.gov.hk/en/tide/predtide.htm` | Predicted tide tables |
| TC Track (GIS) | `https://www.hko.gov.hk/en/wxinfo/currwx/tc_gis.htm` | GIS-based tropical cyclone track viewer |
| TC Strike Probability | `https://www.hko.gov.hk/en/probfcst/tc_spm.htm` | Tropical cyclone strike probability maps |
| UV Index Forecast | `https://www.hko.gov.hk/en/wxinfo/uvinfo/uvinfo.html` | UV index nowcast and forecast |
| Earth Weather Viewer | `https://maps.weather.gov.hk/wxviewer/index.html?lang=en` | Global weather visualization |
| World Weather | `https://www.hko.gov.hk/en/wxinfo/worldwx/wwi.htm` | World city weather information |
| Past Weather | `https://www.hko.gov.hk/en/wxinfo/pastwx/past.htm` | Historical weather data archive |
| King's Park AWS | `https://www.hko.gov.hk/en/wxinfo/aws/kpinfo.htm` | Automatic weather station data |

### 5.5 RSS Polling Strategy for Discord Bot

| RSS Feed | Recommended Interval | Use Case |
|---|---|---|
| `CurrentWeather.xml` | 10 min | Lightweight current conditions |
| `LocalWeatherForecast.xml` | 60 min | Today/tomorrow forecast |
| `SeveralDaysWeatherForecast_v2.xml` | 120 min | Extended forecast |
| `WeatherWarningBulletin.xml` | 1-2 min | Warning alert monitoring |
| `WeatherWarningSummaryv2.xml` | 1-2 min | Quick warning status |
| `QuickEarthquakeMessage.xml` | 5-10 min | Earthquake notifications |
| `FeltEarthquake.xml` | 5-10 min | Local tremor reports |

---

## 6. Rainfall in The Past Hour from Automatic Weather Station API

**Endpoint:** `https://data.weather.gov.hk/weatherAPI/opendata/hourlyRainfall.php`
**Method:** GET
**Return Type:** JSON

**Parameters:** `lang=en|tc|sc`

**Example:** `?lang=en`

```json
{
    "obsTime": "2026-07-26T00:00:00+08:00",
    "hourlyRainfall": [
        {
            "automaticWeatherStation": "Hong Kong Observatory",
            "automaticWeatherStationID": "HKO",
            "value": 0,
            "unit": "mm"
        }
    ]
}
```

> **Note:** This data is provisional and from automatic weather stations. It may differ from the official HKO rainfall record.

---

## 7. Discord Bot Implementation Patterns

### Polling Strategy

| Endpoint | Recommended Interval | Use Case |
|---|---|---|
| `weather.php?dataType=rhrread` | 5-10 min | Current conditions display |
| `weather.php?dataType=fnd` | 30-60 min | 9-day forecast embed |
| `weather.php?dataType=flw` | 30-60 min | Forecast text channel |
| `weather.php?dataType=warnsum` | 1-2 min | Warning alert monitoring |
| `weather.php?dataType=warningInfo` | On warnsum change | Detailed warning text |
| `weather.php?dataType=swt` | 5 min | Special tips |
| `earthquake.php?dataType=qem` | 5-10 min | Earthquake alerts |
| `hourlyRainfall.php` | 5-10 min | Live rainfall |

### Suggested Command Structure

```
/weather current       — rhrread: temp, humidity, UV, rainfall
/weather forecast      — fnd: next 3 days
/weather forecast 9    — fnd: full 9 days
/weather warnings      — warnsum: active warnings
/weather warning <type> — warningInfo: detailed bulletin
/weather tide          — HHOT: today's tides
/weather lunar         — lunardate: today's lunar date
/weather rain          — hourlyRainfall: past hour
/weather earthquake    — qem: latest quakes
/weather uv            — rhrread UV index
```

### Warning Alert Monitoring Pattern

```
loop every 1-2 min:
    fetch warnsum
    compare with previous warnsum state
    if new warning issued or escalated:
        fetch warningInfo for that warning
        send rich embed to alert channel
    if warning cancelled:
        send cancellation notice
```

### Rate Limiting

There is no documented rate limit, but be a good citizen:
- Cache responses client-side where appropriate (especially `fnd` which only updates every 30-60 min)
- Use `ETag`/`Last-Modified` if the server sends them
- Minimum 1-2 seconds between requests
- Avoid polling `rhrread` more than once per minute

### Language Support

All endpoints accept `lang=en|tc|sc`. Consider making the bot bilingual or letting users set their preferred language via a server config command.

---

*End of reference. For full documentation, see [HKO Open Data API Documentation PDF](https://data.weather.gov.hk/weatherAPI/doc/HKO_Open_Data_API_Documentation.pdf).*
