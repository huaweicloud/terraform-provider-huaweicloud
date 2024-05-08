package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRoutingBacklogPolicyResponse Response Object
type ListRoutingBacklogPolicyResponse struct {

	// 数据流转积压策略列表。
	BacklogPolicies *[]BacklogPolicyInfo `json:"backlog_policies,omitempty"`

	// 满足查询条件的记录总数。
	Count *int32 `json:"count,omitempty"`

	// 本次分页查询结果中最后一条记录的ID，可在下一次分页查询时使用。
	Marker         *string `json:"marker,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListRoutingBacklogPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRoutingBacklogPolicyResponse struct{}"
	}

	return strings.Join([]string{"ListRoutingBacklogPolicyResponse", string(data)}, " ")
}
