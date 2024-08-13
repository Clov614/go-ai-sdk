// Package ai_sdk
// @Author Clover
// @Data 2024/8/12 下午8:04:00
// @Desc
package ai_sdk

import (
	"ai-sdk/config"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestAIClient_SendByFuncCall(t *testing.T) {
	type fields struct {
		Authorization string
		ContentType   string
		Model         string
		Url           string
		EndPoint      string
		ProxyAddr     string // *url.URL
		timeout       int
	}
	type args struct {
		req Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResp   Response[FunctionCallResponse]
		isWantResp bool
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "TestFunctionCall",
			fields: fields{
				Authorization: "sk-6m0S1PzHRCxtdrSGbh9ABH5uIdQYHKXqNxg4umhJ2LJ5pypc",
				ContentType:   config.DefaultContentType,
				Model:         config.DefaultModel,
				Url:           "https://api.chatanywhere.tech",
				EndPoint:      config.DefaultEndPoint,
				ProxyAddr:     "127.0.0.1:7890",
				timeout:       0,
			},
			args: args{
				req: Request{
					Messages: []ChatMessage{
						{
							Role:    "user",
							Content: "这是一条测试消息",
						},
					},
					Tools: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAIClient(tt.fields.Authorization, tt.fields.Model, tt.fields.Url, tt.fields.ProxyAddr, tt.fields.EndPoint)
			gotResp, err := a.SendByFuncCall(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendByFuncCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.isWantResp {
				if !reflect.DeepEqual(gotResp, tt.wantResp) {
					t.Errorf("SendByFuncCall() gotResp = %v, want %v", gotResp, tt.wantResp)
				}
			} else {
				t.Logf("SendByFuncCall() gotResp = %v", gotResp)
			}
		})
	}
}

func TestSend(t *testing.T) {
	resp, err := aiclient.Send(Request{
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: "这是一条测试消息",
			},
		},
	})
	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	t.Log(resp)
}

func TestAIClient_Send(t *testing.T) {
	type fields struct {
		Authorization string
		ContentType   string
		Model         string
		Url           string
		ProxyAddr     *url.URL
		client        *http.Client
		timeout       int
	}
	type args struct {
		req Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp Response[DefalutResponse]
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AIClient{
				Authorization: tt.fields.Authorization,
				ContentType:   tt.fields.ContentType,
				Model:         tt.fields.Model,
				Url:           tt.fields.Url,
				ProxyAddr:     tt.fields.ProxyAddr,
				client:        tt.fields.client,
				timeout:       tt.fields.timeout,
			}
			gotResp, err := a.Send(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendByStreaming() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("SendByStreaming() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
