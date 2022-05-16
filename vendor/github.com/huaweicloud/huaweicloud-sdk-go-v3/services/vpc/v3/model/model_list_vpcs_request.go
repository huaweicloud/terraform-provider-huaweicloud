package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListVpcsRequest struct {

	// 功能说明：每页返回的个数 取值范围：0-2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`

	// VPC资源ID。可以使用该字段过滤VPC
	Id *[]string `json:"id,omitempty"`

	// VPC的name信息，可以使用该字段过滤VPC
	Name *[]string `json:"name,omitempty"`

	// VPC的描述信息。可以使用该字段过滤VPC
	Description *[]string `json:"description,omitempty"`

	// VPC的CIDR。可以使用该字段过滤VPC
	Cidr *[]string `json:"cidr,omitempty"`
}

func (o ListVpcsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVpcsRequest struct{}"
	}

	return strings.Join([]string{"ListVpcsRequest", string(data)}, " ")
}
