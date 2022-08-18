package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DomainItemDetail struct {

	// 数据起始时间戳，可能与请求时间不一致
	StartTime *int64 `json:"start_time,omitempty"`

	// 数据结束时间戳，可能与请求时间不一致
	EndTime *int64 `json:"end_time,omitempty"`

	// 指标类型
	StatType *string `json:"stat_type,omitempty"`

	// 指标统计数据列表，如果该时间段内无值，则为空数组[]
	Domains *[]map[string]interface{} `json:"domains,omitempty"`
}

func (o DomainItemDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainItemDetail struct{}"
	}

	return strings.Join([]string{"DomainItemDetail", string(data)}, " ")
}
