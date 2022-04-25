package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowCustomPolicyRequest struct {

	// 待查询的自定义策略ID，获取方式请参见：[自定义策略ID](https://apiexplorer.developer.huaweicloud.com/apiexplorer/doc?product=IAM&api=ListCustomPolicies)。
	RoleId string `json:"role_id"`
}

func (o ShowCustomPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCustomPolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowCustomPolicyRequest", string(data)}, " ")
}
