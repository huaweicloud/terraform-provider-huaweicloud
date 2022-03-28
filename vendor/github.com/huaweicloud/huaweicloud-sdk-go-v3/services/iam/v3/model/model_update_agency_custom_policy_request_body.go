package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateAgencyCustomPolicyRequestBody struct {
	Role *AgencyPolicyRoleOption `json:"role"`
}

func (o UpdateAgencyCustomPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAgencyCustomPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateAgencyCustomPolicyRequestBody", string(data)}, " ")
}
