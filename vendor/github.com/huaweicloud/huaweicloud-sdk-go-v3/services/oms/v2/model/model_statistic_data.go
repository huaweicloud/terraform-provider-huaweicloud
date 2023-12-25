package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StatisticData 查询同步任务统计结果
type StatisticData struct {

	// 统计时间戳
	TimeStamp *int64 `json:"time_stamp,omitempty"`

	// 统计数量
	StatisticNum *int64 `json:"statistic_num,omitempty"`
}

func (o StatisticData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StatisticData struct{}"
	}

	return strings.Join([]string{"StatisticData", string(data)}, " ")
}
