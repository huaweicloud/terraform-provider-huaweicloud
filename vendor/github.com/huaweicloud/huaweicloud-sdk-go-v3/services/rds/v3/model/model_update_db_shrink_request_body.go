package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateDbShrinkRequestBody struct {

	// 数据库名称。
	DbName string `json:"db_name"`
}

func (o UpdateDbShrinkRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDbShrinkRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateDbShrinkRequestBody", string(data)}, " ")
}
