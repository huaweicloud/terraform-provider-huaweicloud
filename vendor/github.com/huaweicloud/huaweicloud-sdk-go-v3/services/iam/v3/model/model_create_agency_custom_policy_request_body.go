package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateAgencyCustomPolicyRequestBody struct {
	Role *AgencyPolicyRoleOption `json:"role"`
}

func (o CreateAgencyCustomPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgencyCustomPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"CreateAgencyCustomPolicyRequestBody", string(data)}, " ")
}
