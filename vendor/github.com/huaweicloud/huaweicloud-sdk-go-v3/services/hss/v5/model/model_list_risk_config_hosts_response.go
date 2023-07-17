package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRiskConfigHostsResponse Response Object
type ListRiskConfigHostsResponse struct {

	// 数据总量
	TotalNum *int64 `json:"total_num,omitempty"`

	// 数据列表
	DataList       *[]SecurityCheckHostInfoResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                                  `json:"-"`
}

func (o ListRiskConfigHostsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRiskConfigHostsResponse struct{}"
	}

	return strings.Join([]string{"ListRiskConfigHostsResponse", string(data)}, " ")
}
