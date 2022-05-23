package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRoutingRulesResponse struct {

	// 规则条件信息列表。
	Rules *[]RoutingRule `json:"rules,omitempty"`

	// 满足查询条件的记录总数。
	Count *int32 `json:"count,omitempty"`

	// 本次分页查询结果中最后一条记录的ID，可在下一次分页查询时使用。
	Marker         *string `json:"marker,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListRoutingRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRoutingRulesResponse struct{}"
	}

	return strings.Join([]string{"ListRoutingRulesResponse", string(data)}, " ")
}
