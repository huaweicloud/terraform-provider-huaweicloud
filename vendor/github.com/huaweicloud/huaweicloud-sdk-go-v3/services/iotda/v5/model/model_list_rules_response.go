package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRulesResponse struct {

	// 本次分页查询结果中最后一条记录的ID，可在下一次分页查询时使用。
	Marker *string `json:"marker,omitempty"`

	// 满足查询条件的记录总数。
	Count *int64 `json:"count,omitempty"`

	// 规则信息列表。
	Rules          *[]RuleResponse `json:"rules,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRulesResponse struct{}"
	}

	return strings.Join([]string{"ListRulesResponse", string(data)}, " ")
}
