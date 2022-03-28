package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type AssociateRoleToUserOnEnterpriseProjectResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AssociateRoleToUserOnEnterpriseProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateRoleToUserOnEnterpriseProjectResponse struct{}"
	}

	return strings.Join([]string{"AssociateRoleToUserOnEnterpriseProjectResponse", string(data)}, " ")
}
