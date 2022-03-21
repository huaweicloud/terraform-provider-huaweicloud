package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateCloudServiceCustomPolicyResponse struct {
	Role           *ServicePolicyRoleResult `json:"role,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o UpdateCloudServiceCustomPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCloudServiceCustomPolicyResponse struct{}"
	}

	return strings.Join([]string{"UpdateCloudServiceCustomPolicyResponse", string(data)}, " ")
}
