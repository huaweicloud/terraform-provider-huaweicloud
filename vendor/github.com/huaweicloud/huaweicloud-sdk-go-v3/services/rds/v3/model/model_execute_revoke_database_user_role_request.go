package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExecuteRevokeDatabaseUserRoleRequest Request Object
type ExecuteRevokeDatabaseUserRoleRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *DatabaseUserRoleRequest `json:"body,omitempty"`
}

func (o ExecuteRevokeDatabaseUserRoleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExecuteRevokeDatabaseUserRoleRequest struct{}"
	}

	return strings.Join([]string{"ExecuteRevokeDatabaseUserRoleRequest", string(data)}, " ")
}
