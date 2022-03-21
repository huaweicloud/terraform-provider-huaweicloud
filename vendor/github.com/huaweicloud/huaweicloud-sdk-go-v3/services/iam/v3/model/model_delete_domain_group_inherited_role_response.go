package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteDomainGroupInheritedRoleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteDomainGroupInheritedRoleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainGroupInheritedRoleResponse struct{}"
	}

	return strings.Join([]string{"DeleteDomainGroupInheritedRoleResponse", string(data)}, " ")
}
