package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleEnterpriseProjectAssignmentId struct {

	// 企业项目ID。
	Id *string `json:"id,omitempty"`
}

func (o RoleEnterpriseProjectAssignmentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleEnterpriseProjectAssignmentId struct{}"
	}

	return strings.Join([]string{"RoleEnterpriseProjectAssignmentId", string(data)}, " ")
}
