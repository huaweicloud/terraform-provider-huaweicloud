package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleDomainAssignmentId struct {

	// 全局服务ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleDomainAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleDomainAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleDomainAssignmentId", string(data)}, " ")
}
