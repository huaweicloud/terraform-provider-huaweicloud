package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListUsersForGroupByAdminRequest struct {

	// 待查询的用户组ID，获取方式请参见：[获取用户组ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	GroupId string `json:"group_id"`
}

func (o KeystoneListUsersForGroupByAdminRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListUsersForGroupByAdminRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListUsersForGroupByAdminRequest", string(data)}, " ")
}
