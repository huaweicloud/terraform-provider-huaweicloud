package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ResetGaussMySqlPasswordRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID，严格匹配UUID规则。

	InstanceId string `json:"instance_id"`

	Body *MysqlResetPasswordRequest `json:"body,omitempty"`
}

func (o ResetGaussMySqlPasswordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetGaussMySqlPasswordRequest struct{}"
	}

	return strings.Join([]string{"ResetGaussMySqlPasswordRequest", string(data)}, " ")
}
