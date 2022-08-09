package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TopUrlSummary struct {

	// URL名称。
	Url *string `json:"url,omitempty"`

	// 对应查询类型的值。（流量单位：Byte）
	Value *int64 `json:"value,omitempty"`

	// 查询起始时间戳。
	StartTime *int64 `json:"start_time,omitempty"`

	// 查询结束时间戳
	EndTime *int64 `json:"end_time,omitempty"`

	// 参数类型支持：flux(流量)，req_num(请求总数)。
	StatType *string `json:"stat_type,omitempty"`
}

func (o TopUrlSummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TopUrlSummary struct{}"
	}

	return strings.Join([]string{"TopUrlSummary", string(data)}, " ")
}
