package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExecuteRevokeDatabaseUserRoleResponse Response Object
type ExecuteRevokeDatabaseUserRoleResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ExecuteRevokeDatabaseUserRoleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExecuteRevokeDatabaseUserRoleResponse struct{}"
	}

	return strings.Join([]string{"ExecuteRevokeDatabaseUserRoleResponse", string(data)}, " ")
}
