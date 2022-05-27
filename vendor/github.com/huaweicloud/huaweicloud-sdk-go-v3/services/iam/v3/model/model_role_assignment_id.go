package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleAssignmentId struct {

	// 权限ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleAssignmentId", string(data)}, " ")
}
