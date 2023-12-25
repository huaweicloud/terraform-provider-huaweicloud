package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IpFrequencyLimitQuery Ip访问限频。
type IpFrequencyLimitQuery struct {

	// 状态，on：打开，off：关闭。
	Status string `json:"status"`

	// 访问阈值，单位：次/秒。
	Qps *int32 `json:"qps,omitempty"`
}

func (o IpFrequencyLimitQuery) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IpFrequencyLimitQuery struct{}"
	}

	return strings.Join([]string{"IpFrequencyLimitQuery", string(data)}, " ")
}
