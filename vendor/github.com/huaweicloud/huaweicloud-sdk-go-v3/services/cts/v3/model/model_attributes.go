package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Attributes 临时安全凭据的属性。
type Attributes struct {

	// 颁发临时安全凭证时的时间（timestamp，为标准UTC时间，毫秒级，13位数字）。
	CreatedAt *string `json:"created_at,omitempty"`

	// 是否已经通过MFA身份认证。
	MfaAuthenticated *string `json:"mfa_authenticated,omitempty"`
}

func (o Attributes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Attributes struct{}"
	}

	return strings.Join([]string{"Attributes", string(data)}, " ")
}
