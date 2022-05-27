package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoContrast struct {

	// 对比度算法名称\"hw-contrast\"。
	Name *string `json:"name,omitempty"`

	// 1 表示视频处理时第一个执行，2表示第二个执行，以此类推；除不执行，各视频处理算法的执行次序不可相同。
	ExecutionOrder *int32 `json:"execution_order,omitempty"`

	// 对比度调节的程度， 值越大， 对比度越高。
	Contrast *string `json:"contrast,omitempty"`

	// 1 表示视频处理时第一个执行，2表示第二个执行，以此类推；除不执行，各视频处理算法的执行次序不可相同。
	Brightness *string `json:"brightness,omitempty"`
}

func (o VideoContrast) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoContrast struct{}"
	}

	return strings.Join([]string{"VideoContrast", string(data)}, " ")
}
