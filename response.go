// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:48:00
// @Desc
package ai_sdk

type RespType int

const (
	IncrementalRespType RespType = iota // 流模式调用
	FuncCallType                        // 回调方法(tools)模式调用
)

type Response[T any | IncrementalResponse | FunctionCallResponse] struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int64  `json:"created"`
	Model    string `json:"model"`
	Data     T
	baseResp BaseResponse
}

type IncrementalResponse struct {
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
}

type FunctionCallResponse struct {
	Choices []FunctionCallChoice `json:"choices"`
	Usage   Usage                `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Delta        Delta       `json:"delta"`
	Logprobs     interface{} `json:"logprobs,omitempty"` // 可以是nil
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type FunctionCallChoice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"` // 可以是nil
	FinishReason string      `json:"finish_reason"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   *string    `json:"content,omitempty"` // Content可能为null
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON格式的字符串
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
