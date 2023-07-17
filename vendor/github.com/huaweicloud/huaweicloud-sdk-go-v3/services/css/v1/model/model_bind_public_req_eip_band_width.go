package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BindPublicReqEipBandWidth 公网带宽信息。
type BindPublicReqEipBandWidth struct {

	// 带宽大小。单位：Mbit/s
	Size int32 `json:"size"`
}

func (o BindPublicReqEipBandWidth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BindPublicReqEipBandWidth struct{}"
	}

	return strings.Join([]string{"BindPublicReqEipBandWidth", string(data)}, " ")
}
