package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateDomainGroupInheritRoleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDomainGroupInheritRoleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainGroupInheritRoleResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainGroupInheritRoleResponse", string(data)}, " ")
}
