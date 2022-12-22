package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListUserStatisticsResponse struct {

	// 账号总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 账户统计信息列表
	DataList       *[]UserStatisticInfoResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o ListUserStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUserStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ListUserStatisticsResponse", string(data)}, " ")
}
