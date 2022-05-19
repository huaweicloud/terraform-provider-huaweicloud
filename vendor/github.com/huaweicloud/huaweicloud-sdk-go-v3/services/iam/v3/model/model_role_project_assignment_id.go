package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleProjectAssignmentId struct {

	// IAM项目ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleProjectAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleProjectAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleProjectAssignmentId", string(data)}, " ")
}
