package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateDatabaseOwnerRequestBody struct {

	// 修改后数据库owner
	Owner string `json:"owner"`

	// 数据库名称
	Database string `json:"database"`
}

func (o UpdateDatabaseOwnerRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDatabaseOwnerRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateDatabaseOwnerRequestBody", string(data)}, " ")
}
