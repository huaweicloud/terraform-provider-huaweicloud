package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UnbindDevicePolicyResponse Response Object
type UnbindDevicePolicyResponse struct {

	// 策略ID。
	PolicyId *string `json:"policy_id,omitempty"`

	// **参数说明**：策略的目标类型。 **取值范围**：device|product|app，device表示设备，product表示产品，app表示整个资源空间。
	TargetType *string `json:"target_type,omitempty"`

	// 成功的目标id列表。
	SuccessTargets *[]string `json:"success_targets,omitempty"`

	// 失败的目标id列表
	FailureTargets *[]DevicePolicyBindOrUnbindFailureDetail `json:"failure_targets,omitempty"`
	HttpStatusCode int                                      `json:"-"`
}

func (o UnbindDevicePolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnbindDevicePolicyResponse struct{}"
	}

	return strings.Join([]string{"UnbindDevicePolicyResponse", string(data)}, " ")
}
