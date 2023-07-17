package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRiskConfigCheckRulesResponse Response Object
type ListRiskConfigCheckRulesResponse struct {

	// 风险总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 数据列表
	DataList       *[]CheckRuleRiskInfoResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o ListRiskConfigCheckRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRiskConfigCheckRulesResponse struct{}"
	}

	return strings.Join([]string{"ListRiskConfigCheckRulesResponse", string(data)}, " ")
}
