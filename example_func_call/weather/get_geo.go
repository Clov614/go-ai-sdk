// Package weather
// @Author Clover
// @Data 2024/8/15 下午5:16:00
// @Desc
package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	getCityCodeErr = errors.New("getCityCode err")
)

// getCityCodeByAddr 根据地址获取城市代码 cityAddress: 城市地址，如: 泉州市永春县 isMultiDay: 是否获取多日天气
func (w *Weather) getCityCodeByAddr(cityAddress string) (cityCode string, err error) {
	// 构建请求参数
	params := url.Values{}
	params.Add("key", w.key)
	params.Add("address", cityAddress)

	// 拼接完整 URL
	reqURL := fmt.Sprintf("%s?%s", geoReqUrl, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("NewWeather NewRequest err: %w", err)
	}
	response, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w Do Request err: %w", getCityCodeErr, err)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w Do Response err: %v", getCityCodeErr, response.Status)
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("%w ReadAll err: %w", getCityCodeErr, err)
	}
	var geoResp GeocodeResponse
	err = json.Unmarshal(respBody, &geoResp)
	if err != nil {
		return "", fmt.Errorf("%w Unmarshal err: %w", getCityCodeErr, err)
	}
	if geoResp.Status != "1" {
		return "", fmt.Errorf("%w GetCityCode err: %v %v", getCityCodeErr, geoResp.Info, geoResp.Infocode)
	}
	return geoResp.Geocodes[0].Adcode, nil
}
