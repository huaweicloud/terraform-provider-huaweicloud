package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDatabaseUserRoleResponse Response Object
type ListDatabaseUserRoleResponse struct {

	// 角色信息
	Roles          *[]string `json:"roles,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListDatabaseUserRoleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDatabaseUserRoleResponse struct{}"
	}

	return strings.Join([]string{"ListDatabaseUserRoleResponse", string(data)}, " ")
}
