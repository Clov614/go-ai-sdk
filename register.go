// Package ai_sdk
// @Author Clover
// @Data 2024/8/12 下午4:00:00
// @Desc 函数调用注册 todo 设计函数调用外部注册
package ai_sdk

import (
	"fmt"
	"strings"
	"sync"
)

const (
	defaultFuncType = "function"
)

var FuncRegister FuncCallRegister

type FuncCallRegister struct {
	Name2Info           map[string]*FuncCallInfo
	keyWord2FuncNameMap // 关键词命中
	filter              // 过滤器
	mu                  sync.RWMutex
}

// Register 注册调用方法
func (fc *FuncCallRegister) Register(finfo *FuncCallInfo, triggerWords []string) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.Name2Info[finfo.Name] = finfo
	for _, triggerWord := range triggerWords {
		if callInfos, ok := fc.keyWord2FuncNameMap[triggerWord]; ok {
			fc.keyWord2FuncNameMap[triggerWord] = append(callInfos, finfo)
		} else {
			fc.keyWord2FuncNameMap[triggerWord] = []*FuncCallInfo{finfo}
		}
	}
	if finfo.CustomTrigger != nil { // 自定义触发器
		fc.filter = append(fc.filter, finfo)
	}
}

// GetCallInfo 根据方法名获取方法调用信息
func (fc *FuncCallRegister) GetCallInfo(name string) *FuncCallInfo {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	callInfo, ok := fc.Name2Info[name]
	if !ok {
		return nil
	}
	return callInfo
}

type FuncInfoNameList []string

// GetToolsByContent 根据触发条件返回调用方法信息
func (fc *FuncCallRegister) GetToolsByContent(content string) *[]Tool {
	fc.mu.RLock()
	var attackFunction []Function = make([]Function, 0)
	// 先找出所有关键词触发的函数
	for keyWord, infos := range fc.keyWord2FuncNameMap {
		if strings.Contains(content, keyWord) {
			attackFunction = append(attackFunction, infos.funcs()...)
		}
	}

	// 自定义触发条件的加入
	for _, info := range fc.filter {
		if info.CustomTrigger(content) {
			attackFunction = append(attackFunction, info.Function)
		}
	}
	// 过滤一遍去重
	existedMap := make(map[string]bool)
	var tempAttackFuncs = make([]Function, 0)
	for _, function := range attackFunction {
		if !existedMap[function.Name] {
			tempAttackFuncs = append(tempAttackFuncs, function)
			existedMap[function.Name] = true
		}
	}
	attackFunction = tempAttackFuncs
	fc.mu.RUnlock()
	// 封装 tools
	var tools = make([]Tool, len(attackFunction))
	for i, function := range attackFunction {
		tools[i] = Tool{
			Type:     defaultFuncType,
			Function: function,
		}
	}
	return &tools
}

// keyWord2FuncNameMap 关键词触发
type keyWord2FuncNameMap map[string]funcInfoList

type funcInfoList []*FuncCallInfo

func (finfos funcInfoList) funcs() []Function {
	var funcs = make([]Function, len(finfos))
	for i, finfo := range finfos {
		funcs[i] = finfo.Function
	}
	return funcs
}

func (finfos funcInfoList) names() FuncInfoNameList {
	funcInfoNameList := make(FuncInfoNameList, len(finfos))
	for i, finfo := range finfos {
		funcInfoNameList[i] = finfo.Name
	}
	return funcInfoNameList
}

// Name2Filter 方法调用过滤器
type filter []*FuncCallInfo

//// FuncCallFilter 方法调用过滤器
//type FuncCallFilter interface {
//	IsCall(content string) bool // 是否调用
//}

type CallFunc interface {
	Call(params string) (jsonStr string, err error) // 调用外部函数，返回json
}

type FuncCallInfo struct {
	Function      // 方法信息
	CallFunc      // 外部函数
	CustomTrigger func(content string) bool
}

func (fc *FuncCallInfo) IsCall(content string) bool {
	if fc.CustomTrigger != nil {
		return fc.CustomTrigger(content)
	}
	return false
}

// Call 调用方法
func (fc *FuncCallInfo) Call(callId string, params string) (Message, error) {
	content, err := fc.CallFunc.Call(params)
	if err != nil {
		return Message{}, fmt.Errorf("FuncCallInfo.CallFunc.Call error: %w", err)
	}
	return Message{
		Role:       toolRole,
		Content:    content,
		ToolCallID: callId,
	}, nil
}

func init() {
	FuncRegister = FuncCallRegister{
		Name2Info:           make(map[string]*FuncCallInfo),
		keyWord2FuncNameMap: make(map[string]funcInfoList),
		filter:              make(filter, 0),
	}
}
