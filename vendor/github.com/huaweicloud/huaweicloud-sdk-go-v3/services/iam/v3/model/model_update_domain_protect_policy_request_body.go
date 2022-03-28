package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateDomainProtectPolicyRequestBody struct {
	ProtectPolicy *ProtectPolicyOption `json:"protect_policy"`
}

func (o UpdateDomainProtectPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainProtectPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateDomainProtectPolicyRequestBody", string(data)}, " ")
}
