package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleAgencyAssignmentId struct {

	// 委托ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleAgencyAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleAgencyAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleAgencyAssignmentId", string(data)}, " ")
}
