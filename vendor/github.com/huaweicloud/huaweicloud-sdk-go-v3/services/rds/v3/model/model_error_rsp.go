package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ErrorRsp 失败时返回的对象。
type ErrorRsp struct {

	// 错误码。
	ErrorCode string `json:"error_code"`

	// 错误描述。
	ErrorMsg string `json:"error_msg"`
}

func (o ErrorRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrorRsp struct{}"
	}

	return strings.Join([]string{"ErrorRsp", string(data)}, " ")
}
