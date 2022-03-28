package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateDomainProtectPolicyResponse struct {
	ProtectPolicy  *ProtectPolicyResult `json:"protect_policy,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o UpdateDomainProtectPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainProtectPolicyResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainProtectPolicyResponse", string(data)}, " ")
}
