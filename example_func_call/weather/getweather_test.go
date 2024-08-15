// Package weather
// @Author Clover
// @Data 2024/8/15 下午5:52:00
// @Desc 天气接口测试
package weather

import (
	"flag"
	ai_sdk "github.com/Clov614/go-ai-sdk"
	"github.com/Clov614/go-ai-sdk/global"
	"testing"
)

var key *string

func init() {
	key = flag.String("apikey", "", "天气api请求key")

}

func TestGetWeather(t *testing.T) {
	flag.Parse()
	weather := NewWeather(*key)
	weatherResp := weather.GetWeatherByCityAddr("福建省泉州市永春县", false)
	if weatherResp.Err != nil {
		t.Error(weatherResp.Err)
	} else {
		t.Log(weatherResp)
	}
	weatherResp2 := weather.GetWeatherByCityAddr("福建省泉州市永春县", true)
	if weatherResp2.Err != nil {
		t.Error(weatherResp2.Err)
	} else {
		t.Log(weatherResp2)
	}
}

func TestRegisterWeatherFunc(t *testing.T) {
	flag.Parse()
	weather := NewWeather(*key)

	funcCallInfo := ai_sdk.FuncCallInfo{
		Function: ai_sdk.Function{
			Name:        "get_weather_by_city",
			Description: "根据地址获取城市代码 cityAddress: 城市地址，如: 泉州市永春县 isMultiDay: 是否获取多日天气",
			Parameters: ai_sdk.FunctionParameter{
				Type: global.ObjType,
				Properties: ai_sdk.Properties{
					"city_addr": ai_sdk.Property{
						Type:        global.StringType,
						Description: "地址，如：国家，城市，县、区地址",
					},
					"is_multi": ai_sdk.Property{
						Type:        global.BoolType,
						Description: "是否获取多日天气",
						Enum:        []string{"true", "false"},
					},
				},
				Required: []string{"city_addr", "is_multi"},
			},
			Strict: false,
		},
		CallFunc: weather,
		//CustomTrigger: nil, // 暂时不测试
	}

	ai_sdk.FuncRegister.Register(&funcCallInfo, []string{"天气", "weather"})

	tools := ai_sdk.FuncRegister.GetToolsByContent("永春最近天气怎么样")
	if *tools == nil || len(*tools) == 0 {
		t.Error(*tools)
	} else {
		t.Log(*tools)
	}
}
