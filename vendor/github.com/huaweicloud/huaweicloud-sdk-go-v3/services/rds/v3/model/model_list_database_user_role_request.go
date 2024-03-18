package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDatabaseUserRoleRequest Request Object
type ListDatabaseUserRoleRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 用户名，有值时返回该账号可以授权的角色集合
	UserName *string `json:"user_name,omitempty"`
}

func (o ListDatabaseUserRoleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDatabaseUserRoleRequest struct{}"
	}

	return strings.Join([]string{"ListDatabaseUserRoleRequest", string(data)}, " ")
}
