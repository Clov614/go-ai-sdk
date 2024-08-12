// Package ai_sdk
// @Author Clover
// @Data 2024/8/12 下午4:13:00
// @Desc 返回错误类
package ai_sdk

type Ret int

const (
	authorizationError  Ret = iota + 1 // 鉴权错误，请检查是否正确填写'OPEN-API-KEY'
	modelUnSupportError                // 模型不支持，请检查配置文件
	proxyUnUsefulError                 // 代理错误，请检查配置文件
)

type BaseResponse struct {
	Ret
	ErrMsg string `json:"err_msg"`
}

func (b BaseResponse) Ok() bool {
	return b.Ret == 0
}

func (b BaseResponse) Err() error {
	if b.Ret == 0 {
		return nil
	}
	return b.Ret
}
