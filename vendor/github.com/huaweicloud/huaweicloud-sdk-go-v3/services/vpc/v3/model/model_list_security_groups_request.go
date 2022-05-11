package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListSecurityGroupsRequest struct {

	// 功能说明：每页返回的个数 取值范围：0-2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`

	// 功能说明：安全组资源ID。可以使用该字段精确过滤安全组，支持多个ID
	Id *[]string `json:"id,omitempty"`

	// 功能说明：安全组名称。可以使用该字段精确过滤满足条件的安全组，支持传入多个name过滤
	Name *[]string `json:"name,omitempty"`

	// 功能说明：安全组描述新增。可以使用该字段精确过滤安全组，支持传入多个描述进行过滤
	Description *[]string `json:"description,omitempty"`

	// 功能说明：企业项目ID。可以使用该字段过滤某个企业项目下的安全组。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。 约束：若需要查询当前用户所有有权限查看企业项目绑定的安全组，请传参all_granted_eps。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListSecurityGroupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSecurityGroupsRequest struct{}"
	}

	return strings.Join([]string{"ListSecurityGroupsRequest", string(data)}, " ")
}
