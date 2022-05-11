package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRolesForGroupOnEnterpriseProjectResponse struct {

	// 角色列表。
	Roles          *[]RolesItem `json:"roles,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListRolesForGroupOnEnterpriseProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRolesForGroupOnEnterpriseProjectResponse struct{}"
	}

	return strings.Join([]string{"ListRolesForGroupOnEnterpriseProjectResponse", string(data)}, " ")
}
