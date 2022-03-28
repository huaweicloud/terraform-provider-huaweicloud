package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateDomainPasswordPolicyResponse struct {
	PasswordPolicy *PasswordPolicyResult `json:"password_policy,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o UpdateDomainPasswordPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainPasswordPolicyResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainPasswordPolicyResponse", string(data)}, " ")
}
