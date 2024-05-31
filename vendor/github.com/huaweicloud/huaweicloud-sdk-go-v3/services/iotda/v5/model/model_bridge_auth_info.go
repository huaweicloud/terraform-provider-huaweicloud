package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BridgeAuthInfo 网桥鉴权信息
type BridgeAuthInfo struct {

	// 鉴权类型。当前支持密钥认证接入(SECRET)。使用密钥认证接入方式(SECRET)填写secret字段，不填写auth_type默认为密钥认证接入方式(SECRET)。
	AuthType *string `json:"auth_type,omitempty"`

	// 网桥密钥，认证类型使用密钥认证接入(SECRET)可填写该字段。
	Secret *string `json:"secret,omitempty"`
}

func (o BridgeAuthInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BridgeAuthInfo struct{}"
	}

	return strings.Join([]string{"BridgeAuthInfo", string(data)}, " ")
}
