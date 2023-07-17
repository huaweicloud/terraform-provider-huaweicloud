package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePublicKibanaBandwidthReqBandWidth 带宽。
type UpdatePublicKibanaBandwidthReqBandWidth struct {

	// 修改后的带宽大小。
	Size int32 `json:"size"`
}

func (o UpdatePublicKibanaBandwidthReqBandWidth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePublicKibanaBandwidthReqBandWidth struct{}"
	}

	return strings.Join([]string{"UpdatePublicKibanaBandwidthReqBandWidth", string(data)}, " ")
}
