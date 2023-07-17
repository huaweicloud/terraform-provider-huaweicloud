package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateAgencyEpPolicyAssignmentReqBody struct {

	// 委托在企业项目上的绑定关系，最多支持250条。
	RoleAssignments []CreateAgencyEpPolicyAssignmentReqBodyRoleAssignments `json:"role_assignments"`
}

func (o CreateAgencyEpPolicyAssignmentReqBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgencyEpPolicyAssignmentReqBody struct{}"
	}

	return strings.Join([]string{"CreateAgencyEpPolicyAssignmentReqBody", string(data)}, " ")
}
