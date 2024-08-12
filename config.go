// Package ai_sdk
// @Author Clover
// @Data 2024/8/12 下午12:53:00
// @Desc 配置项
package ai_sdk

import (
	"ai-sdk/utils/configutil"
	"errors"
	"os"
)

type AICfg struct {
	Authorization string `yaml:"authorization" comment:"OPEN-API-KEY api密钥 (必填)"`
	ContentType   string `yaml:"content_type,omitempty" comment:"默认: application/json"`
	Model         string `yaml:"model,omitempty" comment:"使用的模型ID 默认: gpt-4o-mini"`
	Url           string `yaml:"api_url" comment:"api地址 默认: https://api.openai.com/v1/chat/completions"`
	ProxyAddr     string `yaml:"proxy_address,omitempty" comment:"代理地址 (可选)"`
	// 功能设置项
	Timeout    int `yaml:"timeout" comment:"请求超时时间，单位秒，默认 5s"`
	HistoryNum int `yaml:"history_num,omitempty" comment:"最大上下文长度 默认: 10"`
}

const (
	defaultContentType = "application/json"
	defaultModel       = "gpt-4o-mini"
	defaultUrl         = "https://api.openai.com/v1/chat/completions"
	defaultHistoryNum  = 10 // 默认上下文长度
)

var Config = AICfg{
	ContentType: defaultContentType,
	Model:       defaultModel,
	Url:         defaultUrl,
	HistoryNum:  defaultHistoryNum,
}

var defaultPath = "./cfg/"

var defaultSaveFileName = "ai-cfg.yaml"

func init() {
	err := configutil.Load(&Config, defaultPath, defaultSaveFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = configutil.Save(&Config, defaultPath, defaultSaveFileName)
		}
		ErrorWithErr(err, "error load config")
	}
	//cfg.verifiability() // 校验设置项是否合规
	err = configutil.Save(&Config, defaultPath, defaultSaveFileName)
	if err != nil {
		ErrorWithErr(err, "error saving config")
	}
}
