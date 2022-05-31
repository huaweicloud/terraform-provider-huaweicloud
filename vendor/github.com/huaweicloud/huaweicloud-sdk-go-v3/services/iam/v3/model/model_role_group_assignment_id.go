package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleGroupAssignmentId struct {

	// 用户组ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleGroupAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleGroupAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleGroupAssignmentId", string(data)}, " ")
}
