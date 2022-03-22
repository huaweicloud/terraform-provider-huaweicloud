package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateGaussMySqlBackupPolicyResponse struct {
	// 状态信息

	Status *string `json:"status,omitempty"`
	// 实例ID

	InstanceId *string `json:"instance_id,omitempty"`
	// 实例名称

	InstanceName   *string `json:"instance_name,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateGaussMySqlBackupPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateGaussMySqlBackupPolicyResponse struct{}"
	}

	return strings.Join([]string{"UpdateGaussMySqlBackupPolicyResponse", string(data)}, " ")
}
