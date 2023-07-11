package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUsersRequest Request Object
type ListUsersRequest struct {

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// 账号名称
	UserName *string `json:"user_name,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 是否允许登陆
	LoginPermission *bool `json:"login_permission,omitempty"`

	// 是否有root权限
	RootPermission *bool `json:"root_permission,omitempty"`

	// 用户组
	UserGroup *string `json:"user_group,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListUsersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUsersRequest struct{}"
	}

	return strings.Join([]string{"ListUsersRequest", string(data)}, " ")
}
