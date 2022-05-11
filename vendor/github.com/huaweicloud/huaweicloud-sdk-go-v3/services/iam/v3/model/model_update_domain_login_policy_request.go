package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateDomainLoginPolicyRequest struct {

	// 账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	Body *UpdateDomainLoginPolicyRequestBody `json:"body,omitempty"`
}

func (o UpdateDomainLoginPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainLoginPolicyRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainLoginPolicyRequest", string(data)}, " ")
}
