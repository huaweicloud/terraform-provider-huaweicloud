package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UnlockNodeReadonlyStatusRequestBody 解除只读请求体
type UnlockNodeReadonlyStatusRequestBody struct {

	// Ha保持不再设置节点只读状态的时间，单位为分钟。
	StatusPreservationTime int32 `json:"status_preservation_time"`
}

func (o UnlockNodeReadonlyStatusRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnlockNodeReadonlyStatusRequestBody struct{}"
	}

	return strings.Join([]string{"UnlockNodeReadonlyStatusRequestBody", string(data)}, " ")
}
