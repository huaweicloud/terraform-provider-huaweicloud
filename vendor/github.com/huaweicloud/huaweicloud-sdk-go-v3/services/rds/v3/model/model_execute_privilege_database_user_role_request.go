package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExecutePrivilegeDatabaseUserRoleRequest Request Object
type ExecutePrivilegeDatabaseUserRoleRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *DatabaseUserRoleRequest `json:"body,omitempty"`
}

func (o ExecutePrivilegeDatabaseUserRoleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExecutePrivilegeDatabaseUserRoleRequest struct{}"
	}

	return strings.Join([]string{"ExecutePrivilegeDatabaseUserRoleRequest", string(data)}, " ")
}
