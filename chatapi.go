// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:34:00
// @Desc openapi 客户端
package ai_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Clov614/go-ai-sdk/config"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var aiclient *AIClient

type AIClient struct {
	ContentType string
	Model       string
	ApiCfgList  []config.APIConfig
	client      *http.Client
	timeout     int
	EndPoint    string
}

// NewAIClient 创建一个自定义请求客户端
func NewAIClient(apiCfgList []config.APIConfig, model string, endPoint string, timeout int) *AIClient {
	client := AIClient{
		ContentType: config.DefaultContentType,
		Model:       model,
		ApiCfgList:  apiCfgList,
		EndPoint:    endPoint,
	}
	if timeout < 10 {
		client.timeout = 10
	}
	client.client = &http.Client{
		Timeout: time.Duration(client.timeout) * time.Second,
	}
	return &client
}

// 根据 proxy 字符串 设置请求代理
func (a AIClient) transformProxy(proxyAddr string) {
	if proxyAddr != "" {
		proxy, err := url.Parse(ensurePrefix(proxyAddr))
		if err != nil {
			log.Error().Err(err).Str("proxyAddr", proxyAddr).Msg("invalid proxy URL")
			return
		}
		a.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	} else {
		a.client.Transport = nil // 无代理
	}
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

//func (a AIClient) SendFuncCall(content string, tools *[]Tool) (resp Response[DefalutResponse], err error) {
//
//}

// SendByFuncCall  使用 Send默认调用就可以支持 Function_Call
// deprecated
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

	// 循环重试发送请求
	var req *http.Request
	var resp *http.Response

apiCfgLoop:
	for _, apiCfg := range a.ApiCfgList {

		for _, auth := range apiCfg.AuthList {
			// 设置请求的req
			req, err = http.NewRequest(http.MethodPost, apiCfg.Url+a.EndPoint, bytes.NewBuffer(body))
			if err != nil {
				log.Error().Err(err).Msg("new request failed")
				continue
			}
			// 转换代理并设置
			a.transformProxy(apiCfg.ProxyAddr)
			req.Header.Set("Content-Type", a.ContentType)
			// 根据auth尝试进行请求
			req.Header.Set("Authorization", ensureBearer(auth))
			resp, err = a.client.Do(req) // nolint:bodyclose
			if err != nil {
				log.Error().Err(err).Msg("send ai talk request failed")
				continue
			}
			// 根据状态码处理响应
			if statusCode := resp.StatusCode; statusCode == http.StatusOK {
				err = nil                          // 错误置空
				response.baseResp = BaseResponse{} // 错误置空
				break apiCfgLoop
			}
			// 打印错误信息保存错误
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				err = fmt.Errorf("api: %s, %w", apiCfg.Url, unAuthErr)
				response.baseResp = BaseResponse{
					Ret:    authorizationError,
					ErrMsg: "401 Authorization Required",
				}
			case http.StatusMethodNotAllowed:
				err = fmt.Errorf("api: %s, %w", apiCfg.Url, methodNotAllowedErr)
				log.Error().Err(paramUnSupportError).Fields(map[string]interface{}{
					"request":  request,
					"url":      apiCfg.Url,
					"EndPoint": a.EndPoint,
				}).Msg(resp.Status)
				response.baseResp = BaseResponse{
					Ret:    authorizationError,
					ErrMsg: "405 Not Allowed",
				}
			default:
				// 处理其他未预期的状态码
				// nolint
				err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				log.Error().Err(err).Fields(map[string]interface{}{
					"request":  request,
					"url":      apiCfg.Url,
					"EndPoint": a.EndPoint,
				}).Msg(resp.Status)
				response.baseResp = BaseResponse{
					Ret:    authorizationError,
					ErrMsg: fmt.Sprintf("%d Not Allowed", resp.StatusCode),
				}
			}
			resp.Body.Close()
			resp = nil
		}
	}
	if resp == nil {
		response.baseResp = BaseResponse{
			Ret:    paramUnSupportError,
			ErrMsg: "返回值为空，请检查配置文件设置项是否正确填写",
		}
		return response, fmt.Errorf("response empty err: %w", configErr)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return response, fmt.Errorf("all requests failed %w", err)
	}
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
	if response.ID == "" {
		var respErr RespError
		err := json.Unmarshal(body, &respErr)
		if err != nil {
			return response, fmt.Errorf("doSend json.Unmarshal(body, &respErr): %w", err)
		}
		response.err = respErr
	}
	var data T
	err = json.Unmarshal(body, &data)
	if err != nil {
		return response, fmt.Errorf("doSend json.Unmarshal(body, &data): %w", err)
	}
	response.data = data
	return response, nil
}

func init() {
	aiclient = &AIClient{

		ContentType: config.Config.ContentType,
		Model:       config.Config.Model,
		ApiCfgList:  config.Config.ApiCfgs,
		EndPoint:    config.Config.EndPoint,
	}
	if config.Config.Timeout < 10 {
		aiclient.timeout = 10
	} else {
		aiclient.timeout = config.Config.Timeout
	}
	aiclient.client = &http.Client{
		Timeout: time.Duration(aiclient.timeout) * time.Second,
	}
}
