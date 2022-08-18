package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListGroupsRequest struct {

	// 用户组所属账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId *string `json:"domain_id,omitempty"`

	// 用户组名，长度小于等于128字符，获取方式请参见：[获取用户组名](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Name *string `json:"name,omitempty"`
}

func (o KeystoneListGroupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListGroupsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListGroupsRequest", string(data)}, " ")
}
