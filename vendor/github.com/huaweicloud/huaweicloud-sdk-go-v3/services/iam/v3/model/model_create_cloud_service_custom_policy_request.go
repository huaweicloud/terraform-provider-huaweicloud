package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateCloudServiceCustomPolicyRequest struct {
	Body *CreateCloudServiceCustomPolicyRequestBody `json:"body,omitempty"`
}

func (o CreateCloudServiceCustomPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCloudServiceCustomPolicyRequest struct{}"
	}

	return strings.Join([]string{"CreateCloudServiceCustomPolicyRequest", string(data)}, " ")
}
