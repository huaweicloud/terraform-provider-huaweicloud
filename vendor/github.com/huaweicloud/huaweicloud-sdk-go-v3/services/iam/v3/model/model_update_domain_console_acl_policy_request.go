package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateDomainConsoleAclPolicyRequest struct {

	// 账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	Body *UpdateDomainConsoleAclPolicyRequestBody `json:"body,omitempty"`
}

func (o UpdateDomainConsoleAclPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainConsoleAclPolicyRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainConsoleAclPolicyRequest", string(data)}, " ")
}
