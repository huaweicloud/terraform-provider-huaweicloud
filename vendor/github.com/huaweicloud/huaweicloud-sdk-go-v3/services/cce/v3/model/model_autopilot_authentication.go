package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotAuthentication
type AutopilotAuthentication struct {

	// 集群认证模式。默认取值为“rbac”。
	Mode *string `json:"mode,omitempty"`
}

func (o AutopilotAuthentication) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotAuthentication struct{}"
	}

	return strings.Join([]string{"AutopilotAuthentication", string(data)}, " ")
}
