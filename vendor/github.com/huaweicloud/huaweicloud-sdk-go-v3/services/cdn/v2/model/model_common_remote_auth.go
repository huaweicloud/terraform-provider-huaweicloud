package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CommonRemoteAuth 远程鉴权配置。
type CommonRemoteAuth struct {

	// 是否开启远程鉴权(on：开启，off：关闭)。
	RemoteAuthentication string `json:"remote_authentication"`

	RemoteAuthRules *RemoteAuthRule `json:"remote_auth_rules"`
}

func (o CommonRemoteAuth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CommonRemoteAuth struct{}"
	}

	return strings.Join([]string{"CommonRemoteAuth", string(data)}, " ")
}
