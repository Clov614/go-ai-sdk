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
下面是一个使用 go-ai-sdk 的基本示例：

```go
package main

import (
	"fmt"
	"github.com/Clov614/go-ai-sdk"
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

## 配置
该项目使用配置文件来管理各种设置，包括会话超时时间和历史记录长度。你可以在 config.yaml 文件中自定义这些设置：

```yaml
# OPEN-API-KEY api密钥 (必填)
authorization: sk-6m0xxxxxxx5pxxc
# 默认: application/json
content_type: application/json
# 使用的模型ID 默认: gpt-4o-mini
model: gpt-4o-mini
# api地址 默认: https://api.openai.com/v1/chat/completions
api_url: https://api.openai.com/v1/chat/completions
# 请求节点 默认: /v1/chat/completions
end_point: /v1/chat/completions
# 请求超时时间，单位秒，默认 5s
timeout: 0
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