package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UnlockNodeReadonlyStatusRequest Request Object
type UnlockNodeReadonlyStatusRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *UnlockNodeReadonlyStatusRequestBody `json:"body,omitempty"`
}

func (o UnlockNodeReadonlyStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnlockNodeReadonlyStatusRequest struct{}"
	}

	return strings.Join([]string{"UnlockNodeReadonlyStatusRequest", string(data)}, " ")
}
