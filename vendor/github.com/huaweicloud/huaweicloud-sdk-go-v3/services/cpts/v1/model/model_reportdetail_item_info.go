package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReportdetailItemInfo struct {
	// 自定义事务数据

	CustomTransactions *[]string `json:"customTransactions,omitempty"`
	// aw数据

	DetailDatas *[]DetailDataInfo `json:"detailDatas,omitempty"`

	Performance *PerformanceInfo `json:"performance,omitempty"`
}

func (o ReportdetailItemInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportdetailItemInfo struct{}"
	}

	return strings.Join([]string{"ReportdetailItemInfo", string(data)}, " ")
}
