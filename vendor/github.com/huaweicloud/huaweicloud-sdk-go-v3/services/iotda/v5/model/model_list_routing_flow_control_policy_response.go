package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRoutingFlowControlPolicyResponse Response Object
type ListRoutingFlowControlPolicyResponse struct {

	// 数据流转流控策略列表。
	FlowcontrolPolicies *[]FlowControlPolicyInfo `json:"flowcontrol_policies,omitempty"`

	// 满足查询条件的记录总数。
	Count *int32 `json:"count,omitempty"`

	// 本次分页查询结果中最后一条记录的ID，可在下一次分页查询时使用。
	Marker         *string `json:"marker,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListRoutingFlowControlPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRoutingFlowControlPolicyResponse struct{}"
	}

	return strings.Join([]string{"ListRoutingFlowControlPolicyResponse", string(data)}, " ")
}
