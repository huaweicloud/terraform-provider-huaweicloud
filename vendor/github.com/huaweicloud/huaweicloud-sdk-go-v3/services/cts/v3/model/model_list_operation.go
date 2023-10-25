package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOperation 本次查询全量云服务的操作事件列表。
type ListOperation struct {

	// 事件对应的云服务类型。
	ServiceType *string `json:"service_type,omitempty"`

	// 事件对应的资源类型。
	ResourceType *string `json:"resource_type,omitempty"`

	// 操作事件名称数组。
	OperationList *[]string `json:"operation_list,omitempty"`
}

func (o ListOperation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOperation struct{}"
	}

	return strings.Join([]string{"ListOperation", string(data)}, " ")
}
