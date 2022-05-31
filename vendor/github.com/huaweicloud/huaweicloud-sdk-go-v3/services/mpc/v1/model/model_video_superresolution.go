package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoSuperresolution struct {

	// 超分算法名称\"hw-sr\"。
	Name *string `json:"name,omitempty"`

	// 1 表示视频处理时第一个执行，2表示第二个执行，以此类推；除不执行，各视频处理算法的执行次序不可相同。
	ExecutionOrder *int32 `json:"execution_order,omitempty"`

	// 超分倍数，取值范围是[2,8]，默认2。
	Scale *string `json:"scale,omitempty"`
}

func (o VideoSuperresolution) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoSuperresolution struct{}"
	}

	return strings.Join([]string{"VideoSuperresolution", string(data)}, " ")
}
