// Package weather
// @Author Clover
// @Data 2024/8/15 下午5:40:00
// @Desc 获取天气，并注册方法
package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	allWeather  = "all"
	baseWeahter = "base"
)

var (
	GetWeatherErr = errors.New("get weather error")
)

type Weather struct {
	client *http.Client
	key    string
}

// nolint
type DefaultWeatherResp struct {
	WeatherResponse
	MultiDayWeatherResponse
	Err error
}

func NewWeather(apikey string) *Weather {
	return &Weather{
		client: &http.Client{
			Timeout: time.Second * 20,
		},
		key: apikey,
	}
}

// Call AI调用获取天气
func (w *Weather) Call(params string) (jsonStr string, err error) {
	// 注册调用函数以及触发规则
	type weatherProperties struct {
		CityAddr string `json:"city_addr"` // 城市地址
		IsMulti  bool   `json:"is_multi"`
	}
	var properties weatherProperties
	err = json.Unmarshal([]byte(params), &properties)
	if err != nil {
		return "", fmt.Errorf("call_err: json.Unmarshal([]byte(params)) %w: %w", GetWeatherErr, err)
	}
	resp := w.GetWeatherByCityAddr(properties.CityAddr, properties.IsMulti)
	if resp.Err != nil {
		return "", fmt.Errorf("call_err: GetWeatherByCityAddr() error: %w", resp.Err)
	}
	bytes, err := json.Marshal(resp) // nolint
	if err != nil {
		return "", fmt.Errorf("call_err: %w json.Marshal() error: %w", GetWeatherErr, err)
	}
	return string(bytes), nil
}

// GetWeatherByCityAddr 根据国家，城市，县、区地址获取天气 isMultiDay: 是否获取多日天气
func (w *Weather) GetWeatherByCityAddr(cityAddr string, isMultiDay bool) (weatherResp DefaultWeatherResp) {
	cityCode, err := w.getCityCodeByAddr(cityAddr)
	if err != nil {
		weatherResp.Err = err
		return
	}

	// 构建请求参数
	params := url.Values{}
	params.Add("key", w.key)
	params.Add("city", cityCode)
	if isMultiDay {
		params.Add("extensions", allWeather)
	} else {
		params.Add("extensions", baseWeahter)
	}
	// 拼接完整 URL
	reqURL := fmt.Sprintf("%s?%s", weatherReqUrl, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		weatherResp.Err = fmt.Errorf("NewWeather NewRequest err: %w", err)
		return
	}
	response, err := w.client.Do(req)
	if err != nil {
		weatherResp.Err = fmt.Errorf("%w Do Request err: %w", GetWeatherErr, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		weatherResp.Err = fmt.Errorf("%w Do Response err: %v", GetWeatherErr, response.Status)
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		weatherResp.Err = fmt.Errorf("%w ReadAll err: %w", GetWeatherErr, err)
		return
	}
	if isMultiDay {
		var multiResp MultiDayWeatherResponse
		err = json.Unmarshal(respBody, &multiResp)
		if err != nil {
			weatherResp.Err = fmt.Errorf("%w Unmarshal err: %w", GetWeatherErr, err)
			return
		}
		weatherResp.MultiDayWeatherResponse = multiResp
		return
	}
	var resp WeatherResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		weatherResp.Err = fmt.Errorf("%w Unmarshal err: %w", GetWeatherErr, err)
		return
	}
	weatherResp.WeatherResponse = resp
	return
}
