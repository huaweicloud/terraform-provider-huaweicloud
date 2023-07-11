package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProtectionPolicyResponse Response Object
type UpdateProtectionPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateProtectionPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProtectionPolicyResponse struct{}"
	}

	return strings.Join([]string{"UpdateProtectionPolicyResponse", string(data)}, " ")
}
