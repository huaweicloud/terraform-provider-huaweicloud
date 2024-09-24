package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTargetsInDevicePolicyResponse Response Object
type ShowTargetsInDevicePolicyResponse struct {

	// 策略绑定信息列表。
	Targets *[]PolicyTargetBase `json:"targets,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ShowTargetsInDevicePolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTargetsInDevicePolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowTargetsInDevicePolicyResponse", string(data)}, " ")
}
