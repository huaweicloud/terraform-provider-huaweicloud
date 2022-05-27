package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoSaturation struct {

	// 饱和度算法名称\"“hw-saturation\"。
	Name *string `json:"name,omitempty"`

	// 1 表示视频处理时第一个执行，2表示第二个执行，以此类推；除不执行，各视频处理算法的执行次序不可相同。
	ExecutionOrder *int32 `json:"execution_order,omitempty"`

	// 饱和度调节的程度， 值越大， 饱和度越高。
	Saturation *string `json:"saturation,omitempty"`
}

func (o VideoSaturation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoSaturation struct{}"
	}

	return strings.Join([]string{"VideoSaturation", string(data)}, " ")
}
