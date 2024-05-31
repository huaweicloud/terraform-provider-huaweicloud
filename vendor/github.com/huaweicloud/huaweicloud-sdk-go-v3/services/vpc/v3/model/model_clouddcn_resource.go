package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClouddcnResource 查询过滤标签的资源
type ClouddcnResource struct {

	// 资源ID标识符。
	ResourceId string `json:"resource_id"`

	// 资源详情。
	ResourceDetail *interface{} `json:"resource_detail"`

	// 包含标签。
	Tags []Tag `json:"tags"`

	// 实例名字。
	ResourceName string `json:"resource_name"`
}

func (o ClouddcnResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClouddcnResource struct{}"
	}

	return strings.Join([]string{"ClouddcnResource", string(data)}, " ")
}
