// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:48:00
// @Desc gpt实体类 request
package ai_sdk

type Request struct {
	Messages   []ChatMessage
	Tools      *[]Tool
	ToolChoice string
}

type ChatCompletionRequest struct {
	Model      string        `json:"model"`
	Messages   []ChatMessage `json:"messages"`
	Tools      *[]Tool       `json:"tools,omitempty"`       // 可选
	ToolChoice string        `json:"tool_choice,omitempty"` // 默认 auto
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// region FunctionCall Request

// FunctionParameter 定义函数参数类型
type FunctionParameter struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

// Property 定义函数属性类型
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"` // 用于枚举类型的字段
}

// Tool 定义函数类型的工具
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function 定义函数结构
type Function struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  FunctionParameter `json:"parameters"`
}

//endregion

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}
