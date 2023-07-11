package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociateRoleToAgencyOnEnterpriseProjectResponse Response Object
type AssociateRoleToAgencyOnEnterpriseProjectResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AssociateRoleToAgencyOnEnterpriseProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateRoleToAgencyOnEnterpriseProjectResponse struct{}"
	}

	return strings.Join([]string{"AssociateRoleToAgencyOnEnterpriseProjectResponse", string(data)}, " ")
}
