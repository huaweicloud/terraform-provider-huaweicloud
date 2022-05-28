package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type AddApplicationRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	Body *AddApplication `json:"body,omitempty"`
}

func (o AddApplicationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddApplicationRequest struct{}"
	}

	return strings.Join([]string{"AddApplicationRequest", string(data)}, " ")
}
