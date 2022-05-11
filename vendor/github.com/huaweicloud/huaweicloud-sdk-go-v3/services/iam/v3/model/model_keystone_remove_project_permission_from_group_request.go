package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneRemoveProjectPermissionFromGroupRequest struct {

	// 项目ID，获取方式请参见：[获取项目名称、项目ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	ProjectId string `json:"project_id"`

	// 用户组ID，获取方式请参见：[获取用户组ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	GroupId string `json:"group_id"`

	// 权限ID，获取方式请参见：[获取权限名、权限ID](https://support.huaweicloud.com/api-iam/iam_10_0001.html)。
	RoleId string `json:"role_id"`
}

func (o KeystoneRemoveProjectPermissionFromGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneRemoveProjectPermissionFromGroupRequest struct{}"
	}

	return strings.Join([]string{"KeystoneRemoveProjectPermissionFromGroupRequest", string(data)}, " ")
}
