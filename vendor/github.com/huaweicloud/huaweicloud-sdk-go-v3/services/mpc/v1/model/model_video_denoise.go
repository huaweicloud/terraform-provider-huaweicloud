package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoDenoise struct {

	// 降噪算法名称\"hw-denoise\"、\"waifu2x\"。
	Name *string `json:"name,omitempty"`

	// 1 表示视频处理时第一个执行，2表示第二个执行，以此类推；除不执行，各视频处理算法的执行次序不可相同。
	ExecutionOrder *int32 `json:"execution_order,omitempty"`
}

func (o VideoDenoise) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoDenoise struct{}"
	}

	return strings.Join([]string{"VideoDenoise", string(data)}, " ")
}
