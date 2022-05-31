package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 资源结构体。
type ResourceDto struct {

	// 资源id。例如，要查询的资源类型为device，那么对应的资源id就是device_id。
	ResourceId *string `json:"resource_id,omitempty"`
}

func (o ResourceDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceDto struct{}"
	}

	return strings.Join([]string{"ResourceDto", string(data)}, " ")
}
