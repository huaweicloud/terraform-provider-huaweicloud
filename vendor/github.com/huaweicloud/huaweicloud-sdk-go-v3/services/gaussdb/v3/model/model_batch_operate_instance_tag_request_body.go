package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchOperateInstanceTagRequestBody struct {
	// 操作标识，取值： - create，表示添加标签。 - delete，表示删除标签。

	Action string `json:"action"`
	// 标签列表。

	Tags []TagItem `json:"tags"`
}

func (o BatchOperateInstanceTagRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchOperateInstanceTagRequestBody struct{}"
	}

	return strings.Join([]string{"BatchOperateInstanceTagRequestBody", string(data)}, " ")
}
