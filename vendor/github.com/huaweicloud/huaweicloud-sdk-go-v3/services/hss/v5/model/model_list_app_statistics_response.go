package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppStatisticsResponse Response Object
type ListAppStatisticsResponse struct {

	// 进程统计信息总数,
	TotalNum *int32 `json:"total_num,omitempty"`

	// 进程统计信息列表
	DataList       *[]AppStatisticResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o ListAppStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ListAppStatisticsResponse", string(data)}, " ")
}
