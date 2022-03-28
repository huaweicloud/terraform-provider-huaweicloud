package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RevokeRoleFromUserOnEnterpriseProjectResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RevokeRoleFromUserOnEnterpriseProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeRoleFromUserOnEnterpriseProjectResponse struct{}"
	}

	return strings.Join([]string{"RevokeRoleFromUserOnEnterpriseProjectResponse", string(data)}, " ")
}
