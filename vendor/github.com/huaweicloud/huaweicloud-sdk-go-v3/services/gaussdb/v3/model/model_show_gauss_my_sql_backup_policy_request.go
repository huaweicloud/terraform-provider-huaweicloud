package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlBackupPolicyRequest struct {
	// 语言。

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID，严格匹配UUID规则。

	InstanceId string `json:"instance_id"`
}

func (o ShowGaussMySqlBackupPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlBackupPolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlBackupPolicyRequest", string(data)}, " ")
}
