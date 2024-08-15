// Package weather
// @Author Clover
// @Data 2024/8/15 下午5:10:00
// @Desc
package weather

const (
	geoReqUrl     = "https://restapi.amap.com/v3/geocode/geo"         // 获取geo地址的url
	weatherReqUrl = "https://restapi.amap.com/v3/weather/weatherInfo" // 获取天气请求地址
)

type WeatherRequest struct {
	Key        string `json:"key"`        // 必要 请求api-key
	City       string `json:"city"`       // city_code 通过geo查询得出
	Extensions string `json:"extensions"` // base or all  base: 实时天气，all: 预报
}

// WeatherResponse 表示天气查询接口的响应结构体
type WeatherResponse struct {
	Status   string `json:"status"`   // 请求状态码，1 表示成功
	Count    string `json:"count"`    // 返回结果数
	Info     string `json:"info"`     // 返回状态说明，如 "OK"
	Infocode string `json:"infocode"` // 返回状态码详细说明
	Lives    []Live `json:"lives"`    // 实时天气信息数组
}

// Live 表示单个城市或区域的实时天气信息
type Live struct {
	Province         string `json:"province"`          // 省份名称
	City             string `json:"city"`              // 城市名称或区县名称
	Adcode           string `json:"adcode"`            // 行政区划代码
	Weather          string `json:"weather"`           // 天气情况，如 "大雨"
	Temperature      string `json:"temperature"`       // 当前温度（摄氏度）
	WindDirection    string `json:"winddirection"`     // 风向，如 "东南"
	WindPower        string `json:"windpower"`         // 风力等级，如 "≤3"
	Humidity         string `json:"humidity"`          // 湿度百分比，如 "96"
	ReportTime       string `json:"reporttime"`        // 天气数据的发布时间
	TemperatureFloat string `json:"temperature_float"` // 当前温度的浮点数表示
	HumidityFloat    string `json:"humidity_float"`    // 湿度的浮点数表示
}

// MultiDayWeatherResponse 表示多日天气预报接口的响应结构体
// nolint
type MultiDayWeatherResponse struct {
	Status    string     `json:"status"`    // 请求状态码，1 表示成功
	Count     string     `json:"count"`     // 返回结果数
	Info      string     `json:"info"`      // 返回状态说明，如 "OK"
	Infocode  string     `json:"infocode"`  // 返回状态码详细说明
	Forecasts []Forecast `json:"forecasts"` // 多日天气预报数组
}

// Forecast 表示一个城市或区域的多日天气预报
type Forecast struct {
	City       string `json:"city"`       // 城市名称或区县名称
	Adcode     string `json:"adcode"`     // 行政区划代码
	Province   string `json:"province"`   // 省份名称
	ReportTime string `json:"reporttime"` // 预报发布时间
	Casts      []Cast `json:"casts"`      // 每日天气预报数组
}

// Cast 表示某一天的天气预报
type Cast struct {
	Date           string  `json:"date"`            // 日期，如 "2024-08-15"
	Week           string  `json:"week"`            // 星期几，1-7 分别表示周一到周日
	DayWeather     string  `json:"dayweather"`      // 白天天气情况
	NightWeather   string  `json:"nightweather"`    // 夜间天气情况
	DayTemp        string  `json:"daytemp"`         // 白天温度（摄氏度）
	NightTemp      string  `json:"nighttemp"`       // 夜间温度（摄氏度）
	DayWind        string  `json:"daywind"`         // 白天风向
	NightWind      string  `json:"nightwind"`       // 夜间风向
	DayPower       string  `json:"daypower"`        // 白天风力等级
	NightPower     string  `json:"nightpower"`      // 夜间风力等级
	DayTempFloat   float64 `json:"daytemp_float"`   // 白天温度的浮点数表示
	NightTempFloat float64 `json:"nighttemp_float"` // 夜间温度的浮点数表示
}

// GeocodeResponse 表示地理编码查询接口的响应结构体
// nolint
type GeocodeResponse struct {
	Status   string    `json:"status"`   // 请求状态码，1 表示成功
	Info     string    `json:"info"`     // 返回状态说明，如 "OK"
	Infocode string    `json:"infocode"` // 返回状态码详细说明
	Count    string    `json:"count"`    // 返回结果数
	Geocodes []Geocode `json:"geocodes"` // 地理编码结果数组
}

// Geocode 表示地理编码查询结果
type Geocode struct {
	FormattedAddress string       `json:"formatted_address"` // 格式化后的地址
	Country          string       `json:"country"`           // 国家
	Province         string       `json:"province"`          // 省份
	CityCode         string       `json:"citycode"`          // 城市区号
	City             string       `json:"city"`              // 城市名称
	District         interface{}  `json:"district"`          // 区/县名称
	Township         []string     `json:"township"`          // 乡镇名称，可能为空
	Neighborhood     Neighborhood `json:"neighborhood"`      // 社区信息
	Building         Building     `json:"building"`          // 建筑物信息
	Adcode           string       `json:"adcode"`            // 行政区划代码
	Street           []string     `json:"street"`            // 街道名称，可能为空
	Number           []string     `json:"number"`            // 门牌号码，可能为空
	Location         string       `json:"location"`          // 经度和纬度，格式为 "经度,纬度"
	Level            string       `json:"level"`             // 地址级别，如 "区县"
}

// Neighborhood 表示社区信息
type Neighborhood struct {
	Name []string `json:"name"` // 社区名称，可能为空
	Type []string `json:"type"` // 社区类型，可能为空
}

// Building 表示建筑物信息
type Building struct {
	Name []string `json:"name"` // 建筑物名称，可能为空
	Type []string `json:"type"` // 建筑物类型，可能为空
}
