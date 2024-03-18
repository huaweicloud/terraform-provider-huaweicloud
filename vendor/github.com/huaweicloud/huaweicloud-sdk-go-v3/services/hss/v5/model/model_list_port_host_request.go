package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPortHostRequest Request Object
type ListPortHostRequest struct {

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 端口号
	Port int32 `json:"port"`

	// 端口类型
	Type *string `json:"type,omitempty"`

	// 类别，默认为host，包含如下： - host：主机 - container：容器
	Category *string `json:"category,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListPortHostRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPortHostRequest struct{}"
	}

	return strings.Join([]string{"ListPortHostRequest", string(data)}, " ")
}
