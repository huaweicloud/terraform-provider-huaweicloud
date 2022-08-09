package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 强制跳转。
type ForceRedirectConfig struct {

	// 强制跳转开关（on：打开，off：关闭）。
	Status string `json:"status"`

	// 强制跳转类型（http：强制跳转HTTP，https：强制跳转HTTPS）。
	Type *string `json:"type,omitempty"`
}

func (o ForceRedirectConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ForceRedirectConfig struct{}"
	}

	return strings.Join([]string{"ForceRedirectConfig", string(data)}, " ")
}
