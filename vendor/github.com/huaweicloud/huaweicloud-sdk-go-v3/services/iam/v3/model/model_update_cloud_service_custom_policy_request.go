package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateCloudServiceCustomPolicyRequest struct {

	// 待修改的自定义策略ID，获取方式请参见：[自定义策略ID](https://apiexplorer.developer.huaweicloud.com/apiexplorer/doc?product=IAM&api=ListCustomPolicies)。
	RoleId string `json:"role_id"`

	Body *UpdateCloudServiceCustomPolicyRequestBody `json:"body,omitempty"`
}

func (o UpdateCloudServiceCustomPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCloudServiceCustomPolicyRequest struct{}"
	}

	return strings.Join([]string{"UpdateCloudServiceCustomPolicyRequest", string(data)}, " ")
}
