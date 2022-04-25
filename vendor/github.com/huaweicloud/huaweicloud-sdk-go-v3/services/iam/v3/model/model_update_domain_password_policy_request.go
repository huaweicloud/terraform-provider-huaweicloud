package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateDomainPasswordPolicyRequest struct {

	// 账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	Body *UpdateDomainPasswordPolicyRequestBody `json:"body,omitempty"`
}

func (o UpdateDomainPasswordPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainPasswordPolicyRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainPasswordPolicyRequest", string(data)}, " ")
}
