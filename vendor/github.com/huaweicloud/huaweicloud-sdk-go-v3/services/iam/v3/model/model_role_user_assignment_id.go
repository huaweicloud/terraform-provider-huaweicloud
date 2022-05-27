package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleUserAssignmentId struct {

	// IAM用户ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleUserAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleUserAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleUserAssignmentId", string(data)}, " ")
}
