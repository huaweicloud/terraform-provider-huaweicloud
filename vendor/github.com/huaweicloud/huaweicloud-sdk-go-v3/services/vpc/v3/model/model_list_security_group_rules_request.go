package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListSecurityGroupRulesRequest struct {

	// 功能说明：每页返回个数 取值范围：0-2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`

	// 功能说明：安全组规则ID，支持多个ID过滤
	Id *[]string `json:"id,omitempty"`

	// 功能说明：安全组规则所属安全组ID，支持多个ID过滤
	SecurityGroupId *[]string `json:"security_group_id,omitempty"`

	// 功能说明：安全组规则协议，支持多条过滤
	Protocol *[]string `json:"protocol,omitempty"`

	// 功能说明：安全组规则的描述，支持多个描述同时过滤
	Description *[]string `json:"description,omitempty"`

	// 功能说明：远端安全组ID，支持多ID过滤
	RemoteGroupId *[]string `json:"remote_group_id,omitempty"`

	// 功能说明：安全组规则方向
	Direction *string `json:"direction,omitempty"`

	// 功能说明：安全组规则生效策略
	Action *string `json:"action,omitempty"`
}

func (o ListSecurityGroupRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSecurityGroupRulesRequest struct{}"
	}

	return strings.Join([]string{"ListSecurityGroupRulesRequest", string(data)}, " ")
}
