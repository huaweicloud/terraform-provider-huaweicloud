package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListQuotasDetailResponse Response Object
type ListQuotasDetailResponse struct {

	// 包周期配额数
	PacketCycleNum *int32 `json:"packet_cycle_num,omitempty"`

	// 按需配额数
	OnDemandNum *int32 `json:"on_demand_num,omitempty"`

	// 已使用配额数
	UsedNum *int32 `json:"used_num,omitempty"`

	// 空闲配额数
	IdleNum *int32 `json:"idle_num,omitempty"`

	// 正常配额数
	NormalNum *int32 `json:"normal_num,omitempty"`

	// 过期配额数
	ExpiredNum *int32 `json:"expired_num,omitempty"`

	// 冻结配额数
	FreezeNum *int32 `json:"freeze_num,omitempty"`

	// 配额统计列表
	QuotaStatisticsList *[]QuotaStatisticsResponseInfo `json:"quota_statistics_list,omitempty"`

	// 配额总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 配额列表
	DataList       *[]QuotaResourcesResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                           `json:"-"`
}

func (o ListQuotasDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListQuotasDetailResponse struct{}"
	}

	return strings.Join([]string{"ListQuotasDetailResponse", string(data)}, " ")
}
