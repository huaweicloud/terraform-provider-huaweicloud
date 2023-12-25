package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTrafficMirrorFilterRulesRequest Request Object
type ListTrafficMirrorFilterRulesRequest struct {

	// 使用规则ID过滤或排序
	Id *string `json:"id,omitempty"`

	// 使用规则描述过滤
	Description *string `json:"description,omitempty"`

	// 使用筛选条件ID过滤
	TrafficMirrorFilterId *string `json:"traffic_mirror_filter_id,omitempty"`

	// 使用规则方向过滤
	Direction *string `json:"direction,omitempty"`

	// 使用规则协议过滤
	Protocol *string `json:"protocol,omitempty"`

	// 使用规则源网段过滤
	SourceCidrBlock *string `json:"source_cidr_block,omitempty"`

	// 使用规则目的网段过滤
	DestinationCidrBlock *string `json:"destination_cidr_block,omitempty"`

	// 使用规则源端口范围过滤
	SourcePortRange *string `json:"source_port_range,omitempty"`

	// 使用规则目的端口范围过滤
	DestinationPortRange *string `json:"destination_port_range,omitempty"`

	// 使用规则action过滤
	Action *string `json:"action,omitempty"`

	// 使用规则优先级过滤
	Priority *string `json:"priority,omitempty"`

	// 功能说明：每页返回的个数 取值范围：0-2000
	Limit *int32 `json:"limit,omitempty"`

	// 分页查询起始的资源ID，为空时查询第一页
	Marker *string `json:"marker,omitempty"`
}

func (o ListTrafficMirrorFilterRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrafficMirrorFilterRulesRequest struct{}"
	}

	return strings.Join([]string{"ListTrafficMirrorFilterRulesRequest", string(data)}, " ")
}
