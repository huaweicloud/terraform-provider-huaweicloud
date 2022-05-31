package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoSharp struct {

	// 锐化算法名称\"hw-sharp\"。
	Name *string `json:"name,omitempty"`

	// 1 表示视频处理时第一个执行，2表示第二个执行，以此类推；除不执行，各视频处理算法的执行次序不可相同。
	ExecutionOrder *int32 `json:"execution_order,omitempty"`

	// 锐化的程度， 值越大，锐化越强。
	Amount *string `json:"amount,omitempty"`
}

func (o VideoSharp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoSharp struct{}"
	}

	return strings.Join([]string{"VideoSharp", string(data)}, " ")
}
