// Package ai_sdk
// @Author Clover
// @Data 2024/8/12 下午4:33:00
// @Desc 自定义错误
package ai_sdk

import "errors"

// 错误定义
var (
	networkErr = errors.New("network error") // 网络连接错误
)

func (r Ret) Error() string {
	return r.String()
}
