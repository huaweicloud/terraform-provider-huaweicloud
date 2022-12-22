package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListProcessStatisticsResponse struct {

	// 进程统计信息总数,
	TotalNum *int32 `json:"total_num,omitempty"`

	// 进程统计信息列表
	DataList       *[]ProcessStatisticResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                             `json:"-"`
}

func (o ListProcessStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProcessStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ListProcessStatisticsResponse", string(data)}, " ")
}
