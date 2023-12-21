package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageRiskConfigRulesResponse Response Object
type ListImageRiskConfigRulesResponse struct {

	// 风险总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 数据列表
	DataList       *[]ImageRiskConfigsCheckRulesResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                                       `json:"-"`
}

func (o ListImageRiskConfigRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageRiskConfigRulesResponse struct{}"
	}

	return strings.Join([]string{"ListImageRiskConfigRulesResponse", string(data)}, " ")
}
