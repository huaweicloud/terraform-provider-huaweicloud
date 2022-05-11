package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ScopedTokenIdentity struct {

	// 认证方法，该字段内容为“token”。
	Methods []string `json:"methods"`

	Token *ScopedToken `json:"token"`
}

func (o ScopedTokenIdentity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScopedTokenIdentity struct{}"
	}

	return strings.Join([]string{"ScopedTokenIdentity", string(data)}, " ")
}
