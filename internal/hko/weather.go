package hko

// CurrentWeather is the response from dataType=rhrread.
type CurrentWeather struct {
	Temperature               ReadingGroup    `json:"temperature"`
	Humidity                  ReadingGroup    `json:"humidity"`
	Rainfall                  RainfallGroup   `json:"rainfall"`
	Icon                      interface{}     `json:"icon"`
	IconUpdateTime            string          `json:"iconUpdateTime"`
	UVIndex                   ReadingGroup    `json:"uvindex"`
	SpecialWxTips             interface{}     `json:"specialWxTips"`
	WarningMessage            interface{}     `json:"warningMessage"`
	TCMessage                 interface{}     `json:"tcmessage"`
	MintempFrom00To09         interface{}     `json:"mintempFrom00To09"`
	RainfallFrom00To12        interface{}     `json:"rainfallFrom00To12"`
	RainfallLastMonth         interface{}     `json:"rainfallLastMonth"`
	RainfallJanuaryToLastMonth interface{}    `json:"rainfallJanuaryToLastMonth"`
	UpdateTime                string          `json:"updateTime"`
}

// GetCurrentWeather fetches current weather.
func (c *Client) GetCurrentWeather(lang string) (*CurrentWeather, error) {
	var res CurrentWeather
	if err := c.Get(buildWeatherURL("rhrread", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ForecastDay is a single day in the 9-day forecast.
type ForecastDay struct {
	ForecastDate     string      `json:"forecastDate"`
	Week             string      `json:"week"`
	ForecastWind     string      `json:"forecastWind"`
	ForecastWeather  string      `json:"forecastWeather"`
	ForecastMaxtemp  NumberValue `json:"forecastMaxtemp"`
	ForecastMintemp  NumberValue `json:"forecastMintemp"`
	ForecastMaxrh    NumberValue `json:"forecastMaxrh"`
	ForecastMinrh    NumberValue `json:"forecastMinrh"`
	ForecastIcon     int         `json:"ForecastIcon"`
	PSR              string      `json:"PSR"`
}

// ForecastResponse is the response from dataType=fnd.
type ForecastResponse struct {
	GeneralSituation string        `json:"generalSituation"`
	WeatherForecast  []ForecastDay `json:"weatherForecast"`
	UpdateTime       string        `json:"updateTime"`
	SeaTemp          interface{}   `json:"seaTemp,omitempty"`
	SoilTemp         interface{}   `json:"soilTemp,omitempty"`
}

// GetForecast fetches the 9-day weather forecast.
func (c *Client) GetForecast(lang string) (*ForecastResponse, error) {
	var res ForecastResponse
	if err := c.Get(buildWeatherURL("fnd", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// LocalForecast is the response from dataType=flw.
type LocalForecast struct {
	GeneralSituation string `json:"generalSituation"`
	TCInfo           string `json:"tcInfo"`
	FireDangerWarning string `json:"fireDangerWarning"`
	ForecastPeriod   string `json:"forecastPeriod"`
	ForecastDesc     string `json:"forecastDesc"`
	Outlook          string `json:"outlook"`
	UpdateTime       string `json:"updateTime"`
}

// GetLocalForecast fetches the local weather forecast.
func (c *Client) GetLocalForecast(lang string) (*LocalForecast, error) {
	var res LocalForecast
	if err := c.Get(buildWeatherURL("flw", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// WarningSummaryWarning is a single warning summary.
type WarningSummaryWarning struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	Type       string `json:"type"`
	Subtype    string `json:"subtype,omitempty"`
	ActionCode string `json:"actionCode"`
	IssueTime  string `json:"issueTime"`
	UpdateTime string `json:"updateTime"`
}

// WarningSummary is the response from dataType=warnsum.
type WarningSummary map[string]*WarningSummaryWarning

// GetWarningSummary fetches the warning summary.
func (c *Client) GetWarningSummary(lang string) (WarningSummary, error) {
	var res WarningSummary
	if err := c.Get(buildWeatherURL("warnsum", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return res, nil
}

// WarningDetail is a single warning detail entry.
type WarningDetail struct {
	Contents             []string `json:"contents"`
	WarningStatementCode string   `json:"warningStatementCode"`
	Subtype              string   `json:"subtype"`
	UpdateTime           string   `json:"updateTime"`
}

// WarningInfo is the response from dataType=warningInfo.
type WarningInfo struct {
	Details []WarningDetail `json:"details"`
}

// GetWarningInfo fetches detailed warning information.
func (c *Client) GetWarningInfo(lang string) (*WarningInfo, error) {
	var res WarningInfo
	if err := c.Get(buildWeatherURL("warningInfo", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// SWTTip is a single special weather tip.
type SWTTip struct {
	Desc       string `json:"desc"`
	UpdateTime string `json:"updateTime"`
}

// SWTResponse is the response from dataType=swt.
type SWTResponse struct {
	SWT []SWTTip `json:"swt"`
}

// GetSpecialWeatherTips fetches special weather tips.
func (c *Client) GetSpecialWeatherTips(lang string) (*SWTResponse, error) {
	var res SWTResponse
	if err := c.Get(buildWeatherURL("swt", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// HourlyRainfallResponse is the response from dataType=hourlyRainfall.
type HourlyRainfallResponse struct {
	HourlyRainfall []RainfallGroup `json:"hourlyRainfall"`
}

// GetHourlyRainfall fetches hourly rainfall data.
func (c *Client) GetHourlyRainfall(lang string) (*HourlyRainfallResponse, error) {
	var res HourlyRainfallResponse
	if err := c.Get(buildWeatherURL("hourlyRainfall", languageCode(lang)), &res); err != nil {
		return nil, err
	}
	return &res, nil
}
