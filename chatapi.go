// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:34:00
// @Desc openapi 客户端
package ai_sdk

import (
	"ai-sdk/cfg"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"time"
)

var aiclient *AIClient

type AIClient struct {
	Authorization string
	ContentType   string
	Model         string
	Url           string
	ProxyAddr     *url.URL
	client        *http.Client
	timeout       int
}

func NewAIClient(authorization string, model string, sendurl string, proxyAddr string) *AIClient {
	client := AIClient{
		Authorization: authorization,
		Model:         model,
		Url:           sendurl,
	}
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		log.Error().Err(err).Msg("invalid proxy URL")
	} else {
		client.ProxyAddr = proxy
	}
	if cfg.Config.Timeout < 5 {
		client.timeout = 5
	}
	client.timeout = cfg.Config.Timeout
	client.client = &http.Client{
		Timeout: time.Duration(client.timeout) * time.Second,
	}
	return &client
}

func (a AIClient) SendByStreaming(req Request) (resp Response[IncrementalResponse], err error) {
	resp, err = doSend[IncrementalResponse](a, a.convertReq(req))
	if err != nil {
		return resp, fmt.Errorf("send incremental response failed: %s", err.Error())
	}
	return resp, nil
}

func (a AIClient) SendByFuncCall(req Request) (resp Response[FunctionCallResponse], err error) {
	resp, err = doSend[FunctionCallResponse](a, a.convertReq(req))
	if err != nil {
		return resp, fmt.Errorf("send functioncall response error: %w", err)
	}
	return resp, nil
}

func (a AIClient) convertReq(req Request) ChatCompletionRequest {
	if req.ToolChoice == "" {
		req.ToolChoice = "auto"
	}
	return ChatCompletionRequest{
		Model:         a.Model,
		Messages:      req.Messages,
		Stream:        req.Stream,
		StreamOptions: &StreamOptions{IncludeUsage: req.IncludeUsage},
		Tools:         req.Tools,
		ToolChoice:    req.ToolChoice,
	}
}

func doSend[T IncrementalResponse | FunctionCallResponse](a AIClient, request ChatCompletionRequest) (response Response[T], err error) {
	// 构造请求body
	body, err := json.Marshal(request)
	if err != nil {
		return response, fmt.Errorf("ChatCompletionRequest marshalling failed: %w", err)
	}
	// 发送请求
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest(http.MethodPost, a.Url, bytes.NewBuffer(body))
	if err != nil {
		return response, fmt.Errorf("NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", a.ContentType)
	req.Header.Set("Authorization", a.Authorization)

	resp, err = a.client.Do(req)
	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("%w: %w", networkErr, err)
	}
	defer resp.Body.Close()
	// todo 处理请求成功后的错误
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("doSend ReadAll(resp.body): %w", err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("doSend json.Unmarshal(body, &response): %w", err)
	}
	return response, nil
}

func init() {
	aiclient = &AIClient{
		Authorization: "Bearer " + cfg.Config.Authorization,
		ContentType:   cfg.Config.ContentType,
		Model:         cfg.Config.Model,
		Url:           cfg.Config.Url,
	}
	if cfg.Config.ProxyAddr != "" {
		proxy, err := url.Parse(cfg.Config.ProxyAddr)
		if err != nil {
			log.Error().Err(err).Msg("invalid proxy URL")
		} else {
			aiclient.ProxyAddr = proxy
		}
	}
	if cfg.Config.Timeout < 5 {
		aiclient.timeout = 5
	}
	aiclient.timeout = cfg.Config.Timeout
	aiclient.client = &http.Client{
		Timeout: time.Duration(aiclient.timeout) * time.Second,
	}
	aiclient.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(aiclient.ProxyAddr),
	}
}
