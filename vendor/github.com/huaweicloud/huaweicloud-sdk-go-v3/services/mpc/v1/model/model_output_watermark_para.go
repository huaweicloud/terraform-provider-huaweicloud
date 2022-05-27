package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OutputWatermarkPara struct {

	// 水印时长
	TimeDuration *int32 `json:"time_duration,omitempty"`
}

func (o OutputWatermarkPara) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OutputWatermarkPara struct{}"
	}

	return strings.Join([]string{"OutputWatermarkPara", string(data)}, " ")
}
