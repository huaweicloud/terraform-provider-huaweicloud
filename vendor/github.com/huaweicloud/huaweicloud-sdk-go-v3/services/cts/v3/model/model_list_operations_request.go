package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOperationsRequest Request Object
type ListOperationsRequest struct {

	// 事件对应的云服务类型。
	ServiceType *string `json:"service_type,omitempty"`

	// 事件对应的资源类型。传入该参数时，service_type必选。
	ResourceType *string `json:"resource_type,omitempty"`
}

func (o ListOperationsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOperationsRequest struct{}"
	}

	return strings.Join([]string{"ListOperationsRequest", string(data)}, " ")
}
