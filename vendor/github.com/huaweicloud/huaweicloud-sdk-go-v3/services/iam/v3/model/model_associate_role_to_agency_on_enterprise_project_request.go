package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociateRoleToAgencyOnEnterpriseProjectRequest Request Object
type AssociateRoleToAgencyOnEnterpriseProjectRequest struct {
	Body *CreateAgencyEpPolicyAssignmentReqBody `json:"body,omitempty"`
}

func (o AssociateRoleToAgencyOnEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateRoleToAgencyOnEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"AssociateRoleToAgencyOnEnterpriseProjectRequest", string(data)}, " ")
}
