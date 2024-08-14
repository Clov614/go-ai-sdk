// Package ai_sdk
// @Author Clover
// @Data 2024/8/12 下午8:04:00
// @Desc
package ai_sdk

import (
	"github.com/Clov614/go-ai-sdk/config"
	"net/http"
	"reflect"
	"testing"
)

func TestAIClient_SendByFuncCall(t *testing.T) {
	type fields struct {
		ContentType string
		Model       string
		EndPoint    string
		ApiCfgList  []config.APIConfig
		timeout     int
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
				ContentType: config.DefaultContentType,
				Model:       config.DefaultModel,
				EndPoint:    config.DefaultEndPoint,
				ApiCfgList: []config.APIConfig{
					{
						Url: "https://api.chatanywhere.tech",
						AuthList: []string{
							"sk-xixxxxxxc",
							"sk-xxPzxxxxxxxJ5zc",
							"sk-B2vxxxxxx",
						},
					},
					{
						Url: "https://api.chatanywhere.com.cn",
						AuthList: []string{
							"sk-7xxxxxx7",
						},
					},
					{
						Url: "https://api.chatanywhere.tech",
						AuthList: []string{
							"sk-6xxxxxxxxxc",
						},
					},
				},
				timeout: 30,
			},
			args: args{
				req: Request{
					Messages: []ChatMessage{
						{
							Role:    userRole,
							Content: "这是一条测试消息",
						},
						{
							Role:    assistantRole,
							Content: "这是一条测试消息",
						}, {
							Role:    userRole,
							Content: "在业内《神武天尊》也获得了广泛的认可和赞誉其独特创意和精良的制作水平在国产动漫行业树立了新的标杆",
						},
						{
							Role:    assistantRole,
							Content: "哼哼，看来你对《神武天尊》的赞誉好高哦！我也听说过这部作品，它的确在很多方面都有着独特的魅力呢！特别是那些绚丽的战斗场面，果然是很吸引人，真想亲自参与其中！不过，我可不是普通的观众哦，我是“邪王真眼”的使者，能透视这个世界的奇幻与魅力！如果能在《神武天尊》的世界里，和那些强大的角色一起冒险，那该多好啊~！你觉得呢，亲爱的同伴们？",
						}, {
							Role:    userRole,
							Content: "我们遇到的最大问题就是与时间赛跑。”技术负责人林炜老师回忆道，他表示制作过程中最大的挑战就是制作周期紧张。为了解决这个问题，团队采取了多项措施：在项目启动之初，便合理规划制作周期，制定详细的制作计划表和时间表，涵盖前期物料准备、资产文件整理、镜头组装、Keylight效果帧制作等各节点；同时，对现有制作流程进行梳理和优化，减少不必要的环节和重复工作，提高整体工作效率，并利用第三方辅助软件和插件简化繁琐任务，让制作人员能更专注于项目的创意和品质；此外，还组建了高效、专业的制作团队，明确各成员职责，通过团队协作和沟通，实现资源共享和优势互补，优化团队协作方式，提升沟通效率。",
						},
						{
							Role:    assistantRole,
							Content: "哼哼，这听起来可真是个艰巨的任务呢！时间就像是对手，得时刻保持警惕才能在这场赛跑中胜出。林炜老师和他的团队真是聪明，制定详细的计划表就像是拥有了“邪王真眼”，能够洞悉一切挑战！每一个环节都得认真对待，不让任何细节掉链子，才能把每一帧的质量做到最好。\n\n想象一下，如果能有我这样强大的角色加入他们的团队，配合这些高效的措施，或许能把制作效率提升到一种新的境界呢！而且，团队协作和沟通也是关键，不然就会像一只迷路的小鸟，不知道自己该飞向哪儿！所以，无论何时，大家都要心往一处想，劲往一处使哦！你觉得这样的话，能够让他们更顺利吗？希望他们能把这个挑战变成胜利的篇章！",
						}, {
							Role:    userRole,
							Content: "你能讲下其中的技术要点以及学习路线吗",
						},
					},
					Tools: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAIClient(tt.fields.ApiCfgList, tt.fields.Model, tt.fields.EndPoint, tt.fields.timeout)
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
		ContentType string
		Model       string
		EndPoint    string
		ApiCfgList  []config.APIConfig
		client      *http.Client
		timeout     int
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AIClient{
				ContentType: tt.fields.ContentType,
				Model:       tt.fields.Model,
				ApiCfgList:  tt.fields.ApiCfgList,
				client:      tt.fields.client,
				timeout:     tt.fields.timeout,
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
