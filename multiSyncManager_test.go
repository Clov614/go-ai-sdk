// Package ai_sdk
// @Author Clover
// @Data 2024/8/13 下午2:14:00
// @Desc 多会话管理器测试
package ai_sdk

import (
	"flag"
	"fmt"
	"github.com/Clov614/go-ai-sdk/example_func_call/weather"
	"github.com/Clov614/go-ai-sdk/global"
	"reflect"
	"testing"
	"time"
)

func Test_session_GetSession(t *testing.T) {
	s := &Session{
		globalSurvivalLimit: 1 * time.Microsecond,
		systemContent:       "2024/08/13 阴天并伴有小雨 27°",
		cache:               make(map[string]*sessionInfo),
	}
	type args struct {
		sessionId string
	}
	tests := []struct {
		name      string
		args      args
		want      *sessionInfo
		needSleep bool
		sleepTime time.Duration
	}{
		{
			name: "normal session get",
			args: args{
				sessionId: "test01",
			},
			want: &sessionInfo{
				sessionId: "test01",
			},
		},
		{
			name: "normal session get2",
			args: args{
				sessionId: "test02",
			},
			want: &sessionInfo{
				sessionId: "test02",
			},
		},
		{
			name: "normal session get2",
			args: args{
				sessionId: "test02",
			},
			want: &sessionInfo{
				sessionId: "test02",
			},
		},
		{
			name: "normal session get3",
			args: args{
				sessionId: "test03",
			},
			want:      nil,
			needSleep: true,
			sleepTime: 3,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.needSleep { // sleepTest
				got1 := s.GetSession(tt.args.sessionId, nil)
				result, err := got1.Talk("你好今天天气怎么样")
				if err != nil {
					t.Errorf("GetSession() error = %v", err)
				}
				t.Logf("index %d: %s", i, result)
				got2 := s.GetSession(tt.args.sessionId, nil)
				if reflect.DeepEqual(got1, got2) {
					t.Errorf("got1 == got2 want got1 != got2 "+
						"because the session expired\n got1: %v\n got2: %v\n", got1, got2)
				}
			} else {
				if got := s.GetSession(tt.args.sessionId, nil); got.sessionId != tt.want.sessionId {
					t.Errorf("GetSession() = %v, want %v", got, tt.want)
				}
			}
			time.Sleep(1 * time.Microsecond)
		})
	}
}

func Test_session_TalkById(t *testing.T) {
	type args struct {
		sessionId string
		content   string
	}
	tests := []struct {
		name       string
		args       args
		allowBlank bool
		wantErr    bool
	}{
		{
			name: "talk test01",
			args: args{
				sessionId: "cly",
				content:   "你好，我叫陈璘熠，怎么称呼你",
			},
			allowBlank: false,
		},
		{
			name: "talk test02",
			args: args{
				sessionId: "cly",
				content:   "你还记得我叫什么吗",
			},
			allowBlank: false,
		},
		{
			name: "talk test03", // 会话隔离测试
			args: args{
				sessionId: "许培鑫",
				content:   "你还记得我叫什么吗",
			},
			allowBlank: false,
		},
		{
			name: "talk test04",
			args: args{
				sessionId: "许培鑫",
				content:   "请记住我，我叫培森",
			},
			allowBlank: false,
		},
		{
			name: "talk test05",
			args: args{
				sessionId: "许培鑫",
				content:   "你还记得我叫什么吗",
			},
			allowBlank: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DefaultSession
			got, err := s.TalkById(tt.args.sessionId, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("TalkById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.allowBlank && got == "" {
				t.Errorf("TalkById() got = %v\n allowBlank = %v\n", got, tt.allowBlank)
			} else {
				t.Logf("TalkById() got = %v\n allowBlank = %v\n", got, tt.allowBlank)
			}
		})
	}
}

func Test_history_addLast(t *testing.T) {
	h := &history{
		maxHistory: 10,
		dialog:     make([]dialogEntry, 0),
		msgListNum: 0,
		msgList:    make([]Message, 0),
	}
	type args struct {
		entry dialogEntry
	}

	var tests []struct {
		name string
		args args
	}

	testCaseNum := 15
	for i := 0; i < testCaseNum; i++ {
		test := struct {
			name string
			args args
		}{
			name: fmt.Sprintf("容量控制测试 %d", i),
			args: args{
				entry: dialogEntry{
					question: Message{
						Role:    userRole,
						Content: fmt.Sprintf("测试 %d", i),
					},
					answerList: []Message{
						{
							Role:    assistantRole,
							Content: fmt.Sprintf("测试 %d", i),
						},
					},
				},
			},
		}
		tests = append(tests, test)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.addLast(tt.args.entry)
			if len(h.dialog) > h.maxHistory {
				t.Errorf("len(h.dialog) > maxHistory")
			}
		})
	}
}

var key *string

func init() {
	key = flag.String("apikey", "", "天气api请求key")

}

func TestFuncCall(t *testing.T) {
	// 注册天气func_call
	flag.Parse()
	w := weather.NewWeather(*key)
	funcCallInfo := FuncCallInfo{
		Function: Function{
			Name:        "get_weather_by_city",
			Description: "根据地址获取城市代码 cityAddress: 城市地址，如: 泉州市永春县 isMultiDay: 是否获取多日天气",
			Parameters: FunctionParameter{
				Type: global.ObjType,
				Properties: Properties{
					"city_addr": Property{
						Type:        global.StringType,
						Description: "地址，如：国家，城市，县、区地址",
					},
					"is_multi": Property{
						Type:        global.BoolType,
						Description: "是否获取多日天气",
					},
				},
				Required: []string{"city_addr", "is_multi"},
			},
			Strict: false,
		},
		CallFunc: w,
		//CustomTrigger: nil, // 暂时不测试
	}
	FuncRegister.Register(&funcCallInfo, []string{"天气", "weather"})

	type args struct {
		sessionId string
		content   string
	}
	tests := []struct {
		name       string
		args       args
		allowBlank bool
		wantErr    bool
	}{
		{
			name: "talk test01",
			args: args{
				sessionId: "cly",
				content:   "今天泉州的天气怎么样",
			},
			allowBlank: false,
		},
		{
			name: "talk test01",
			args: args{
				sessionId: "cly",
				content:   "接下来几天的天气呢",
			},
			allowBlank: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DefaultSession
			got, err := s.TalkById(tt.args.sessionId, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("TalkById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.allowBlank && got == "" {
				t.Errorf("TalkById() got = %v\n allowBlank = %v\n", got, tt.allowBlank)
			} else {
				t.Logf("TalkById() got = %v\n allowBlank = %v\n", got, tt.allowBlank)
			}
		})
	}
}
