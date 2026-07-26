package i18n

// M is a map of message keys to translations.
var M = map[string]StringSet{
	"current_weather": {
		EN: "Current Weather",
		TC: "現時天氣",
		SC: "现时天气",
	},
	"forecast": {
		EN: "Weather Forecast",
		TC: "天氣預報",
		SC: "天气预报",
	},
	"forecast_9day": {
		EN: "9-Day Weather Forecast",
		TC: "9天天氣預報",
		SC: "9天天气预报",
	},
	"warnings": {
		EN: "Weather Warnings",
		TC: "天氣警告",
		SC: "天气警告",
	},
	"no_active_warnings": {
		EN: "No active warnings",
		TC: "沒有生效警告",
		SC: "没有生效警告",
	},
	"warning_detail": {
		EN: "Warning Detail",
		TC: "警告詳情",
		SC: "警告详情",
	},
	"rainfall": {
		EN: "Hourly Rainfall",
		TC: "過去一小時雨量",
		SC: "过去一小时雨量",
	},
	"uv_index": {
		EN: "UV Index",
		TC: "紫外線指數",
		SC: "紫外线指数",
	},
	"tides": {
		EN: "Tide Heights",
		TC: "潮汐高度",
		SC: "潮汐高度",
	},
	"lunar_calendar": {
		EN: "Lunar Calendar",
		TC: "農曆",
		SC: "农历",
	},
	"earthquake": {
		EN: "Earthquake Report",
		TC: "地震報告",
		SC: "地震报告",
	},
	"no_earthquake": {
		EN: "No significant earthquake reports",
		TC: "沒有重要地震報告",
		SC: "没有重要地震报告",
	},
	"temperature": {
		EN: "Temperature",
		TC: "氣溫",
		SC: "气温",
	},
	"humidity": {
		EN: "Humidity",
		TC: "相對濕度",
		SC: "相对湿度",
	},
	"rain": {
		EN: "Rainfall",
		TC: "雨量",
		SC: "雨量",
	},
	"station": {
		EN: "Station",
		TC: "站點",
		SC: "站点",
	},
	"updated_at": {
		EN: "Updated",
		TC: "更新時間",
		SC: "更新时间",
	},
	"general_situation": {
		EN: "General Situation",
		TC: "天氣概況",
		SC: "天气概况",
	},
	"outlook": {
		EN: "Outlook",
		TC: "展望",
		SC: "展望",
	},
	"day": {
		EN: "Day",
		TC: "日期",
		SC: "日期",
	},
	"weather": {
		EN: "Weather",
		TC: "天氣",
		SC: "天气",
	},
	"wind": {
		EN: "Wind",
		TC: "風向/風速",
		SC: "风向/风速",
	},
	"max_temp": {
		EN: "Max Temp",
		TC: "最高溫度",
		SC: "最高温度",
	},
	"min_temp": {
		EN: "Min Temp",
		TC: "最低溫度",
		SC: "最低温度",
	},
	"humidity_range": {
		EN: "Humidity",
		TC: "相對濕度",
		SC: "相对湿度",
	},
	"psr": {
		EN: "PSR",
		TC: "顯著降雨概率",
		SC: "显著降雨概率",
	},
	"special_tips": {
		EN: "Special Weather Tips",
		TC: "特別天氣提示",
		SC: "特别天气提示",
	},
	"no_special_tips": {
		EN: "No special weather tips",
		TC: "沒有特別天氣提示",
		SC: "没有特别天气提示",
	},
	"settings": {
		EN: "Bot Settings",
		TC: "機械人設定",
		SC: "机器人设定",
	},
	"language": {
		EN: "Language",
		TC: "語言",
		SC: "语言",
	},
	"alert_channel": {
		EN: "Alert Channel",
		TC: "提示頻道",
		SC: "提示频道",
	},
	"alert_channel_disabled": {
		EN: "Disabled",
		TC: "已停用",
		SC: "已停用",
	},
	"bot_status": {
		EN: "Bot Status",
		TC: "機械人狀態",
		SC: "机器人状态",
	},
	"enabled": {
		EN: "Enabled",
		TC: "已啟用",
		SC: "已启用",
	},
	"disabled": {
		EN: "Disabled",
		TC: "已停用",
		SC: "已停用",
	},
	"setup_permission_required": {
		EN: "You need the Manage Server permission to use setup commands.",
		TC: "你需要管理伺服器權限才能使用設定指令。",
		SC: "你需要管理服务器权限才能使用设定指令。",
	},
	"language_set": {
		EN: "Language set to",
		TC: "語言已設定為",
		SC: "语言已设定为",
	},
	"alert_channel_set": {
		EN: "Alert channel set to",
		TC: "提示頻道已設定為",
		SC: "提示频道已设定为",
	},
	"alert_channel_removed": {
		EN: "Alert channel removed.",
		TC: "已移除提示頻道。",
		SC: "已移除提示频道。",
	},
	"tide_station_set": {
		EN: "Default tide station set to",
		TC: "預設潮汐站已設定為",
		SC: "默认潮汐站已设定为",
	},
	"invalid_tide_station": {
		EN: "Invalid tide station code. Use a valid 3-letter station code.",
		TC: "無效的潮汐站代碼。請使用有效的 3 字母代碼。",
		SC: "无效的潮汐站代码。请使用有效的 3 字母代码。",
	},
	"status_set": {
		EN: "Bot status set to",
		TC: "機械人狀態已設定為",
		SC: "机器人状态已设定为",
	},
	"error_fetching_data": {
		EN: "Error fetching data from HKO.",
		TC: "從天文台獲取資料時出錯。",
		SC: "从天文台获取资料时出错。",
	},
	"no_data": {
		EN: "No data available.",
		TC: "沒有資料。",
		SC: "没有资料。",
	},
	"magnitude": {
		EN: "Magnitude",
		TC: "震級",
		SC: "震级",
	},
	"region": {
		EN: "Region",
		TC: "區域",
		SC: "区域",
	},
	"time": {
		EN: "Time",
		TC: "時間",
		SC: "时间",
	},
	"warning_issued": {
		EN: "Warning Issued",
		TC: "已發出警告",
		SC: "已发出警告",
	},
	"warning_updated": {
		EN: "Warning Updated",
		TC: "警告已更新",
		SC: "警告已更新",
	},
	"warning_cancelled": {
		EN: "Warning Cancelled",
		TC: "警告已取消",
		SC: "警告已取消",
	},
}

// T looks up a translation by key.
func T(key string, lang Language) string {
	if s, ok := M[key]; ok {
		return s.T(lang)
	}
	return key
}
