# go-ai-sdk

![](https://img.shields.io/github/go-mod/go-version/Clov614/go-ai-sdk "语言")
![](https://img.shields.io/github/stars/Clov614/go-ai-sdk?style=flat&color=yellow)
[![](https://img.shields.io/github/actions/workflow/status/Clov614/go-ai-sdk/golangci-lint.yml?branch=main)](https://github.com/Clov614/go-ai-sdk/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/license/Clov614/go-ai-sdk)](https://github.com/Clov614/go-ai-sdk/blob/main/LICENSE "许可协议")

go-ai-sdk 是一个使用 Go 语言编写的 SDK，旨在管理基于 AI 的多会话交互，如聊天机器人或自动化客服系统。该包支持并发会话管理、消息历史处理和基于 AI 的响应生成。

## 功能特点
- 多会话管理：创建并管理多个具有唯一会话 ID 的聊天会话。

- 消息历史处理：存储和管理聊天历史，并支持可配置的历史记录限制。

- AI 响应生成：与 AI 客户端接口，基于聊天历史生成智能响应。

## 安装

在安装 go-ai-sdk 包之前，请确保你已经安装了 `Go 1.21`。然后运行以下命令来安装该包：
```shell
go get github.com/Clov614/go-ai-sdk
```

## 使用方法

### 基本示例

```go
package main

import (
	"fmt"
	ai_sdk "github.com/Clov614/go-ai-sdk"
)

func main() {
	sessionID := "example-session"
	content := "你好，你能帮我做什么？"

	// 获取一个会话或创建一个新会话
	session := ai_sdk.DefaultSession.GetSession(sessionID)

	// 发起对话
	response, err := session.Talk(content)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Println("AI 响应:", response)
}
```

### 使用插件扩展 AI 功能
以下代码展示了如何使用函数注册器将自定义功能（如查询天气）注册到 SDK 中：

```go
package ai

import (
	"encoding/json"
	ai_sdk "github.com/Clov614/go-ai-sdk"
	"github.com/Clov614/go-ai-sdk/example_func_call/weather"
	"github.com/Clov614/go-ai-sdk/global"
	"wechat-demo/rikkabot/config"
	"wechat-demo/rikkabot/logging"
)

type weatherCfg struct {
	Key string `json:"key"`
}

func init() {
	// 注册对话插件
	cfg := config.GetConfig()
	wCfgInterface, ok := cfg.GetCustomPluginCfg("weather_ai")
	if !ok {
		cfg.SetCustomPluginCfg("weather_ai", weatherCfg{Key: ""})
		_ = cfg.Update() // 更新设置
		logging.Fatal("weather_ai plugin config loaded empty. Please write the weather API key in config.yaml", 12)
	}
	bytes, _ := json.Marshal(wCfgInterface)
	var wcfg weatherCfg
	json.Unmarshal(bytes, &wcfg)
	w := weather.NewWeather(wcfg.Key)
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
					},
				},
				Required: []string{"city_addr", "is_multi"},
			},
			Strict: false,
		},
		CallFunc: w,
	}
	ai_sdk.FuncRegister.Register(&funcCallInfo, []string{"天气", "weather"})
}
```
在上述代码中，关键词如“天气”或“weather”会自动触发工具函数调用，为 AI 提供额外的能力。

## 配置
该项目使用配置文件来管理各种设置，包括会话超时时间和历史记录长度。你可以在 config.yaml 文件中自定义这些设置：

```yaml
# 默认: application/json
content_type: application/json
# 使用的模型ID 默认: gpt-4o-mini
model: gpt-4o-mini
# 请求节点 默认: /v1/chat/completions
end_point: /v1/chat/completions
# API 配置列表
configs:
  - # api地址 默认: https://api.openai.com/v1/chat/completions
    api_url: https://api.openai.com/v1/chat/completions
    # OPEN-API-KEY api密钥列表 (必填)
    authorization_list:
      - sk-xxxxxx
      - sk-xxxxxx
# 请求超时时间，单位秒，默认 10s
timeout: 30
# 最大上下文长度 默认: 10
history_num: 10
# 对话会话超时时间 单位: 分钟 默认: 2 minute
session_time_out: 2

```

## 测试
你可以运行项目中提供的测试用例，前提是配置好`OPEN-API-KEY`：

```shell
go test ./...
```

## 贡献
欢迎贡献！请 fork 此仓库，进行修改后提交 pull request。

## 许可证
此项目基于 AGPL-3.0 许可证进行授权 - 详细信息请参阅 [LICENSE](./LICENSE) 文件。

## 联系方式
如果你有任何问题或建议，欢迎联系我。