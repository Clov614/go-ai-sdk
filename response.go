// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:48:00
// @Desc
package ai_sdk

type Response[T any | DefalutResponse | FunctionCallResponse] struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int64  `json:"created"`
	Model    string `json:"model"`
	data     T
	baseResp BaseResponse
	err      RespError
}

func (r Response[T]) GetData() T {
	return r.data
}

type DefalutResponse struct {
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type FunctionCallResponse struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage,omitempty"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs,omitempty"` // 可以是nil
	FinishReason string      `json:"finish_reason"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"` // Content可能为null
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
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

type RespError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    int    `json:"code"`
}
