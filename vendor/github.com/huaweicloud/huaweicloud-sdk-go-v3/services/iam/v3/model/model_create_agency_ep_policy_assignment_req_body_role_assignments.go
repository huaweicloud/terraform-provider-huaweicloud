package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateAgencyEpPolicyAssignmentReqBodyRoleAssignments struct {

	// 委托id
	AgencyId string `json:"agency_id"`

	// 企业项目id
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 策略id
	RoleId string `json:"role_id"`
}

func (o CreateAgencyEpPolicyAssignmentReqBodyRoleAssignments) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgencyEpPolicyAssignmentReqBodyRoleAssignments struct{}"
	}

	return strings.Join([]string{"CreateAgencyEpPolicyAssignmentReqBodyRoleAssignments", string(data)}, " ")
}
