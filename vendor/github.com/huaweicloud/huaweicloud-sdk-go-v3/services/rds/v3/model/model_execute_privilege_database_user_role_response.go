package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExecutePrivilegeDatabaseUserRoleResponse Response Object
type ExecutePrivilegeDatabaseUserRoleResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ExecutePrivilegeDatabaseUserRoleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExecutePrivilegeDatabaseUserRoleResponse struct{}"
	}

	return strings.Join([]string{"ExecutePrivilegeDatabaseUserRoleResponse", string(data)}, " ")
}
