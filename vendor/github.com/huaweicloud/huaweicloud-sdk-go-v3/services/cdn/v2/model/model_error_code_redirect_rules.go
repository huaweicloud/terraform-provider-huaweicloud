package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ErrorCodeRedirectRules 自定义错误页面。
type ErrorCodeRedirectRules struct {

	// 重定向的错误码，当前支持以下状态码 4xx:400, 403, 404, 405, 414, 416, 451 5xx:500, 501, 502, 503, 504
	ErrorCode int32 `json:"error_code"`

	// 重定向状态码，取值为301或302
	TargetCode int32 `json:"target_code"`

	// 重定向的目标链接
	TargetLink string `json:"target_link"`
}

func (o ErrorCodeRedirectRules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrorCodeRedirectRules struct{}"
	}

	return strings.Join([]string{"ErrorCodeRedirectRules", string(data)}, " ")
}
