package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneDeleteGroupRequest struct {

	// 待删除的用户组ID，获取方式请参见：[获取用户组ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	GroupId string `json:"group_id"`
}

func (o KeystoneDeleteGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneDeleteGroupRequest struct{}"
	}

	return strings.Join([]string{"KeystoneDeleteGroupRequest", string(data)}, " ")
}
