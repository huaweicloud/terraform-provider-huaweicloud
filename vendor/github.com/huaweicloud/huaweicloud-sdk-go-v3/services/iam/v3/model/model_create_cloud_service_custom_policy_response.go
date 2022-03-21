package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateCloudServiceCustomPolicyResponse struct {
	Role           *ServicePolicyRoleResult `json:"role,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o CreateCloudServiceCustomPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCloudServiceCustomPolicyResponse struct{}"
	}

	return strings.Join([]string{"CreateCloudServiceCustomPolicyResponse", string(data)}, " ")
}
