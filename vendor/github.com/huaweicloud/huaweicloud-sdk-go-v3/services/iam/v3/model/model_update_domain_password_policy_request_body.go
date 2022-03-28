package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateDomainPasswordPolicyRequestBody struct {
	PasswordPolicy *PasswordPolicyOption `json:"password_policy"`
}

func (o UpdateDomainPasswordPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainPasswordPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateDomainPasswordPolicyRequestBody", string(data)}, " ")
}
