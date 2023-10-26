package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateInstanceRequestBody struct {

	// 是否迁移数据。 - \"true\"：迁移数据。 - \"false\"：不迁移数据。
	MigrateData *string `json:"migrate_data,omitempty"`
}

func (o UpdateInstanceRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateInstanceRequestBody", string(data)}, " ")
}
