package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type Credential struct {

	// AK/SK和securitytoken的过期时间。
	ExpiresAt string `json:"expires_at"`

	// 获取的AK。
	Access string `json:"access"`

	// 获取的SK。
	Secret string `json:"secret"`

	// securitytoken是将所获的AK、SK等信息进行加密后的字符串。
	Securitytoken string `json:"securitytoken"`
}

func (o Credential) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Credential struct{}"
	}

	return strings.Join([]string{"Credential", string(data)}, " ")
}
