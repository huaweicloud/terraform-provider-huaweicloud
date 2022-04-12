package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RespcodeBrokens struct {
	// 校验失败

	CheckPointFailed *[]float64 `json:"checkPointFailed,omitempty"`
	// 异常请求

	Error *[]float64 `json:"error,omitempty"`
	// 其他失败

	OthersFailed *[]float64 `json:"othersFailed,omitempty"`
	// 解析失败

	ParsedFailed *[]float64 `json:"parsedFailed,omitempty"`
	// 连接被拒

	RefusedFailed *[]float64 `json:"refusedFailed,omitempty"`
	// 成功请求

	Success *[]float64 `json:"success,omitempty"`
	// 超时失败

	Timeout *[]float64 `json:"timeout,omitempty"`
}

func (o RespcodeBrokens) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RespcodeBrokens struct{}"
	}

	return strings.Join([]string{"RespcodeBrokens", string(data)}, " ")
}
