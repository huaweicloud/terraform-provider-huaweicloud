package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCustomPolicyResponse Response Object
type ShowCustomPolicyResponse struct {
	Role           *ShowPolicyRoleResult `json:"role,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ShowCustomPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCustomPolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowCustomPolicyResponse", string(data)}, " ")
}
