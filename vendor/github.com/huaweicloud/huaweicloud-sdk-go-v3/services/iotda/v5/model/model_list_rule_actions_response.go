package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRuleActionsResponse struct {

	// 规则动作信息列表。
	Actions *[]RoutingRuleAction `json:"actions,omitempty"`

	// 满足查询条件的记录总数。
	Count *int32 `json:"count,omitempty"`

	// 本次分页查询结果中最后一条记录的ID，可在下一次分页查询时使用。
	Marker         *string `json:"marker,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListRuleActionsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRuleActionsResponse struct{}"
	}

	return strings.Join([]string{"ListRuleActionsResponse", string(data)}, " ")
}
