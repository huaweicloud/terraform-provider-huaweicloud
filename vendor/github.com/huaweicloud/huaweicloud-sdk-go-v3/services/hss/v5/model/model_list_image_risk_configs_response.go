package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageRiskConfigsResponse Response Object
type ListImageRiskConfigsResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 配置检测列表
	DataList       *[]ImageRiskConfigsInfoResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                                 `json:"-"`
}

func (o ListImageRiskConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageRiskConfigsResponse struct{}"
	}

	return strings.Join([]string{"ListImageRiskConfigsResponse", string(data)}, " ")
}
