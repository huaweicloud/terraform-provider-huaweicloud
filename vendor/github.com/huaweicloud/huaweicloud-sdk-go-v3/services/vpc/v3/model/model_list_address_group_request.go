package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAddressGroupRequest Request Object
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

	// 功能说明：企业项目ID。可以使用该字段过滤某个企业项目下的IP地址组。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。若需要查询当前用户所有企业项目绑定的IP地址组，请传参all_granted_eps。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListAddressGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAddressGroupRequest struct{}"
	}

	return strings.Join([]string{"ListAddressGroupRequest", string(data)}, " ")
}
