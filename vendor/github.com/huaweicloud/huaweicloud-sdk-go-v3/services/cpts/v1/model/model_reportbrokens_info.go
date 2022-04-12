package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReportbrokensInfo struct {
	BrandBrokens *BrandBrokens `json:"brand_brokens,omitempty"`
	// 时间戳

	CommonTimestamps *[]string `json:"commonTimestamps,omitempty"`

	RespcodeBrokens *RespcodeBrokens `json:"respcode_brokens,omitempty"`

	TpsBrokens *TpsBrokens `json:"tps_brokens,omitempty"`

	VusersBrokens *VusersBrokens `json:"vusers_brokens,omitempty"`
}

func (o ReportbrokensInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportbrokensInfo struct{}"
	}

	return strings.Join([]string{"ReportbrokensInfo", string(data)}, " ")
}
