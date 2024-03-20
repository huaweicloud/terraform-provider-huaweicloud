package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBandwidthCalcResponse Response Object
type ShowBandwidthCalcResponse struct {

	// 95峰值，日峰值月平均线信息
	BandwidthCalc  map[string]interface{} `json:"bandwidth_calc,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ShowBandwidthCalcResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBandwidthCalcResponse struct{}"
	}

	return strings.Join([]string{"ShowBandwidthCalcResponse", string(data)}, " ")
}
