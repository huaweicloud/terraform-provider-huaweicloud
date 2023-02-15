package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListUserChangeHistoriesRequest struct {

	// 账号名
	UserName *string `json:"user_name,omitempty"`

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 是否有root权限
	RootPermission *bool `json:"root_permission,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 变更类型:   - ADD ：添加   - DELETE ：删除   - MODIFY ： 修改
	ChangeType *string `json:"change_type,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 变更开始时间，13位时间戳
	StartTime *int64 `json:"start_time,omitempty"`

	// 变更结束时间，13位时间戳
	EndTime *int64 `json:"end_time,omitempty"`
}

func (o ListUserChangeHistoriesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUserChangeHistoriesRequest struct{}"
	}

	return strings.Join([]string{"ListUserChangeHistoriesRequest", string(data)}, " ")
}
