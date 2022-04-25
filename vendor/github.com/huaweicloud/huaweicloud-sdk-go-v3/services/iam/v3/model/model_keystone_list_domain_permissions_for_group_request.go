package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListDomainPermissionsForGroupRequest struct {

	// 用户组所属账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	// 用户组ID，获取方式请参见：[获取用户组ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	GroupId string `json:"group_id"`
}

func (o KeystoneListDomainPermissionsForGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListDomainPermissionsForGroupRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListDomainPermissionsForGroupRequest", string(data)}, " ")
}
