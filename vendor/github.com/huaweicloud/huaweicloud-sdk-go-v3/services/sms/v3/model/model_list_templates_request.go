package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTemplatesRequest struct {
	// 模板名称

	Name *string `json:"name,omitempty"`
	// 可用区

	AvailabilityZone *string `json:"availability_zone,omitempty"`
	// Region ID

	Region *string `json:"region,omitempty"`
	// 分页大小，不传值默认为50

	Limit *int32 `json:"limit,omitempty"`
	// 偏移量，不传值默认为0

	Offset *int32 `json:"offset,omitempty"`
}

func (o ListTemplatesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplatesRequest struct{}"
	}

	return strings.Join([]string{"ListTemplatesRequest", string(data)}, " ")
}
