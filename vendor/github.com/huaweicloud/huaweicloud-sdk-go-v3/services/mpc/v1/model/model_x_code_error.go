package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type XCodeError struct {

	// 错误码
	Code *string `json:"code,omitempty"`

	// 错误信息
	Msg *string `json:"msg,omitempty"`
}

func (o XCodeError) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "XCodeError struct{}"
	}

	return strings.Join([]string{"XCodeError", string(data)}, " ")
}
