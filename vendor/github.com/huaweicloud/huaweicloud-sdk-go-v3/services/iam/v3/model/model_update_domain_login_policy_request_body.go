package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateDomainLoginPolicyRequestBody struct {
	LoginPolicy *LoginPolicyOption `json:"login_policy"`
}

func (o UpdateDomainLoginPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainLoginPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateDomainLoginPolicyRequestBody", string(data)}, " ")
}
