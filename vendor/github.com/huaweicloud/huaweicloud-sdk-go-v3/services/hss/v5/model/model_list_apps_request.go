package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListAppsRequest struct {

	// 主机id
	HostId string `json:"host_id"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 版本号
	Version *string `json:"version,omitempty"`

	// 安装目录
	InstallDir *string `json:"install_dir,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListAppsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppsRequest struct{}"
	}

	return strings.Join([]string{"ListAppsRequest", string(data)}, " ")
}
