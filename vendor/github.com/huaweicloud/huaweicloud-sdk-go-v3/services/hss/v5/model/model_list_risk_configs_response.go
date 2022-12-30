package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRiskConfigsResponse struct {

	// 记录总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 服务器配置检测结果列表
	DataList       *[]SecurityCheckInfoResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o ListRiskConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRiskConfigsResponse struct{}"
	}

	return strings.Join([]string{"ListRiskConfigsResponse", string(data)}, " ")
}
