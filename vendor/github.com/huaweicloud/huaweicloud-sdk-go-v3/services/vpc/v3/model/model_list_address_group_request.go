package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListAddressGroupRequest struct {

	// 功能说明：每页返回的个数 取值范围：0~2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`

	// 地址组唯一标识，填写后接口按照id进行过滤，支持多ID同时过滤
	Id *[]string `json:"id,omitempty"`

	// 地址组名称，填写后按照名称进行过滤，支持多名称同时过滤
	Name *[]string `json:"name,omitempty"`

	// IP地址组ip版本，当前只支持ipv4，填写后按照ip版本进行过滤
	IpVersion *int32 `json:"ip_version,omitempty"`

	// 地址组描述信息，填写后按照地址组描述信息过滤，支持多描述同时过滤
	Description *[]string `json:"description,omitempty"`
}

func (o ListAddressGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAddressGroupRequest struct{}"
	}

	return strings.Join([]string{"ListAddressGroupRequest", string(data)}, " ")
}
