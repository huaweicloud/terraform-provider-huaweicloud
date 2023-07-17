package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokeRoleFromAgencyOnEnterpriseProjectRequest Request Object
type RevokeRoleFromAgencyOnEnterpriseProjectRequest struct {
	Body *CreateOrDelAgencyEpPolicyAssignmentReqBody `json:"body,omitempty"`
}

func (o RevokeRoleFromAgencyOnEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeRoleFromAgencyOnEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"RevokeRoleFromAgencyOnEnterpriseProjectRequest", string(data)}, " ")
}
