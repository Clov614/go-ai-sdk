// Code generated by "stringer -type=Ret -linecomment=true -output=stringer.go"; DO NOT EDIT.

package ai_sdk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[authorizationError-1]
	_ = x[modelUnSupportError-2]
	_ = x[proxyUnUsefulError-3]
	_ = x[paramUnSupportError-4]
}

const _Ret_name = "鉴权错误，请检查是否正确填写'OPEN-API-KEY'模型不支持，请检查配置文件代理错误，请检查配置文件参数错误"

var _Ret_index = [...]uint8{0, 56, 95, 131, 143}

func (i Ret) String() string {
	i -= 1
	if i < 0 || i >= Ret(len(_Ret_index)-1) {
		return "Ret(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Ret_name[_Ret_index[i]:_Ret_index[i+1]]
}
