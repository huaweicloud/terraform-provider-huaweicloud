package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 公网ip
type PublicIp struct {
	// 弹性公网IP类型，默认为5_bgp

	Type string `json:"type"`
	// 带宽大小，单位：Mbit/s  调整带宽时的最小单位会根据带宽范围不同存在差异。  小于等于300Mbit/s，默认最小单位为1Mbit/s。300Mbit/s~1000Mbit/s，默认最小单位为50Mbit/s。大于1000Mbit/s：默认最小单位为500Mbit/s。

	BandwidthSize int32 `json:"bandwidth_size"`
}

func (o PublicIp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PublicIp struct{}"
	}

	return strings.Join([]string{"PublicIp", string(data)}, " ")
}
