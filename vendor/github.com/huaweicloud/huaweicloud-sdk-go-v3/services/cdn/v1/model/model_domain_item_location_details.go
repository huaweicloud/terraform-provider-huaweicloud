package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DomainItemLocationDetails struct {

	// 数据起始时间戳，可能与请求时间不一致
	StartTime *int64 `json:"start_time,omitempty"`

	// 数据结束时间戳，可能与请求时间不一致
	EndTime *int64 `json:"end_time,omitempty"`

	// 指标类型
	StatType *string `json:"stat_type,omitempty"`

	// 域名详情数据列表
	Domains *[]DomainRegion `json:"domains,omitempty"`
}

func (o DomainItemLocationDetails) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainItemLocationDetails struct{}"
	}

	return strings.Join([]string{"DomainItemLocationDetails", string(data)}, " ")
}
