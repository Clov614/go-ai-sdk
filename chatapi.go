// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:34:00
// @Desc openapi 客户端
package ai_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	EndPoint      string
}

func NewAIClient(authorization string, model string, sendurl string, proxyAddr string, endPoint string) *AIClient {
	client := AIClient{
		Authorization: ensureBearer(authorization),
		ContentType:   defaultContentType,
		Model:         model,
		Url:           sendurl,
		EndPoint:      endPoint,
	}
	if Config.Timeout < 5 {
		client.timeout = 5
	}
	client.timeout = Config.Timeout
	client.client = &http.Client{
		Timeout: time.Duration(client.timeout) * time.Second,
	}
	if proxyAddr != "" {
		proxy, err := url.Parse(ensurePrefix(proxyAddr))
		if err != nil {
			log.Error().Err(err).Msg("invalid proxy URL")
		} else {
			client.ProxyAddr = proxy
			client.client.Transport = &http.Transport{
				Proxy: http.ProxyURL(aiclient.ProxyAddr),
			}
		}
	}

	return &client
}

// ensureBearer 检查并添加Bearer前缀
func ensureBearer(auth string) string {
	if strings.HasPrefix(auth, "Bearer ") {
		// 如果地址已经有前缀，直接返回
		return auth
	}
	// 如果没有前缀，添加 http://
	return "Bearer " + auth
}

// ensurePrefix 检查并添加http前缀
func ensurePrefix(address string) string {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		// 如果地址已经有前缀，直接返回
		return address
	}
	// 如果没有前缀，添加 http://
	return "http://" + address
}

func (a AIClient) Send(req Request) (resp Response[DefalutResponse], err error) {
	resp, err = doSend[DefalutResponse](a, a.convertReq(req))
	if err != nil {
		return resp, fmt.Errorf("send incremental response failed: %w", err)
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
	if req.Tools != nil && req.ToolChoice == "" {
		req.ToolChoice = "auto"
	}
	return ChatCompletionRequest{
		Model:      a.Model,
		Messages:   req.Messages,
		Tools:      req.Tools,
		ToolChoice: req.ToolChoice,
	}
}

func doSend[T DefalutResponse | FunctionCallResponse](a AIClient, request ChatCompletionRequest) (response Response[T], err error) {
	// 构造请求body
	body, err := json.Marshal(request)
	if err != nil {
		return response, fmt.Errorf("ChatCompletionRequest marshalling failed: %w", err)
	}

	// 发送请求
	req, err := http.NewRequest(http.MethodPost, a.Url+a.EndPoint, bytes.NewBuffer(body))
	if err != nil {
		return response, fmt.Errorf("NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", a.ContentType)
	req.Header.Set("Authorization", a.Authorization)

	resp, err := a.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("%w: %w", networkErr, err)
	}
	defer resp.Body.Close()
	// 根据状态码处理响应
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		response.baseResp = BaseResponse{
			Ret:    authorizationError,
			ErrMsg: "401 Authorization Required",
		}
	case http.StatusMethodNotAllowed:
		log.Error().Err(paramUnSupportError).Fields(map[string]interface{}{
			"request":  request,
			"url":      a.Url,
			"EndPoint": a.EndPoint,
		}).Msg(resp.Status)
		response.baseResp = BaseResponse{
			Ret:    authorizationError,
			ErrMsg: "405 Not Allowed",
		}
	case http.StatusOK:
		// 正常处理响应
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return response, fmt.Errorf("doSend ReadAll(resp.body): %w", err)
		}
		// nolint
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("doSend json.Unmarshal(body, &response): %w", err)
		}
		var data T
		err = json.Unmarshal(body, &data)
		if err != nil {
			return response, fmt.Errorf("doSend json.Unmarshal(body, &data): %w", err)
		}
		if response.ID == "" {
			var respErr RespError
			err := json.Unmarshal(body, &respErr)
			if err != nil {
				return response, fmt.Errorf("doSend json.Unmarshal(body, &respErr): %w", err)
			}
			response.err = respErr
		}
		response.data = data
	default:
		// 处理其他未预期的状态码
		return response, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response, nil
}

func init() {
	aiclient = &AIClient{
		Authorization: "Bearer " + Config.Authorization,
		ContentType:   Config.ContentType,
		Model:         Config.Model,
		Url:           Config.Url,
		EndPoint:      Config.EndPoint,
	}
	if Config.Timeout < 5 {
		aiclient.timeout = 5
	}
	aiclient.timeout = Config.Timeout
	aiclient.client = &http.Client{
		Timeout: time.Duration(aiclient.timeout) * time.Second,
	}
	if Config.ProxyAddr != "" {
		proxy, err := url.Parse(Config.ProxyAddr)
		if err != nil {
			log.Error().Err(err).Msg("invalid proxy URL")
		} else {
			aiclient.ProxyAddr = proxy
			aiclient.client.Transport = &http.Transport{
				Proxy: http.ProxyURL(aiclient.ProxyAddr),
			}
		}
	}
}
