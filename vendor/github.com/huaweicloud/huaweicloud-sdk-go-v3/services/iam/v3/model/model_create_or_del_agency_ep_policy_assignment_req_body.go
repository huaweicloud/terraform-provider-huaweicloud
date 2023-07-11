package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateOrDelAgencyEpPolicyAssignmentReqBody struct {

	// 委托在企业项目上的绑定关系，最多支持250条。
	RoleAssignments []CreateAgencyEpPolicyAssignmentReqBodyRoleAssignments `json:"role_assignments"`
}

func (o CreateOrDelAgencyEpPolicyAssignmentReqBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOrDelAgencyEpPolicyAssignmentReqBody struct{}"
	}

	return strings.Join([]string{"CreateOrDelAgencyEpPolicyAssignmentReqBody", string(data)}, " ")
}
