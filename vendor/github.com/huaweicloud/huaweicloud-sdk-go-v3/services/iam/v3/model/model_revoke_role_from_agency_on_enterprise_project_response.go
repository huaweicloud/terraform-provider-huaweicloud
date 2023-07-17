package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokeRoleFromAgencyOnEnterpriseProjectResponse Response Object
type RevokeRoleFromAgencyOnEnterpriseProjectResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RevokeRoleFromAgencyOnEnterpriseProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeRoleFromAgencyOnEnterpriseProjectResponse struct{}"
	}

	return strings.Join([]string{"RevokeRoleFromAgencyOnEnterpriseProjectResponse", string(data)}, " ")
}
