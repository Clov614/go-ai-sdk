// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:48:00
// @Desc gpt实体类 request
package ai_sdk

type Request struct {
	Messages   []Message
	Tools      *[]Tool
	ToolChoice string
}

type ChatCompletionRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Tools      *[]Tool   `json:"tools,omitempty"`       // 可选
	ToolChoice string    `json:"tool_choice,omitempty"` // 默认 auto
}

const (
	ToolsCallFinishReason = "tool_calls" // 方法调用
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`      // Content可能为null
	ToolCallID string     `json:"tool_call_id,omitempty"` // 用于关联工具调用
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

// region FunctionCall Request

// FunctionParameter 定义函数参数类型
type FunctionParameter struct {
	Type       string `json:"type"`
	Properties `json:"properties"`
	Required   []string `json:"required"`
}

type Properties map[string]Property

// Property 定义函数属性类型
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"` // 用于枚举类型的字段
}

// Tool 定义函数类型的工具
type Tool struct {
	Type     string   `json:"type"` // 默认function
	Function Function `json:"function"`
}

// Function 定义函数结构
type Function struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  FunctionParameter `json:"parameters"`
	Strict      bool              `json:"strict"` // 是否严格 JSON 输出
}

//endregion

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}
