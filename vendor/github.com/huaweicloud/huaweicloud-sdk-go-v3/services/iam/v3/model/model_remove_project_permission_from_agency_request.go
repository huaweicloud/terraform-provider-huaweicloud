package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type RemoveProjectPermissionFromAgencyRequest struct {

	// 委托方的项目ID，获取方式请参见：[获取项目名称、项目ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	ProjectId string `json:"project_id"`

	// 委托ID，获取方式请参见：[获取委托名、委托ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	AgencyId string `json:"agency_id"`

	// 权限ID，获取方式请参见：[获取权限名、权限ID](https://support.huaweicloud.com/api-iam/iam_10_0001.html)。
	RoleId string `json:"role_id"`
}

func (o RemoveProjectPermissionFromAgencyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveProjectPermissionFromAgencyRequest struct{}"
	}

	return strings.Join([]string{"RemoveProjectPermissionFromAgencyRequest", string(data)}, " ")
}
