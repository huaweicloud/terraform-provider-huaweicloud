package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutoLaunchStatisticsResponse Response Object
type ListAutoLaunchStatisticsResponse struct {

	// 自启动项统计信息总数,
	TotalNum *int32 `json:"total_num,omitempty"`

	// 自启动项统计信息列表
	DataList       *[]AutoLaunchStatisticsResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                                 `json:"-"`
}

func (o ListAutoLaunchStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutoLaunchStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ListAutoLaunchStatisticsResponse", string(data)}, " ")
}
