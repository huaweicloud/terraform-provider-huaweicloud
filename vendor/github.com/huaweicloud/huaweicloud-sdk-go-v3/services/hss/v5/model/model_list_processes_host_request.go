package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProcessesHostRequest Request Object
type ListProcessesHostRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 进程可执行文件路径
	Path *string `json:"path,omitempty"`

	// 类型，默认为host，包含如下： - host：主机 - container：容器
	Category *string `json:"category,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListProcessesHostRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProcessesHostRequest struct{}"
	}

	return strings.Join([]string{"ListProcessesHostRequest", string(data)}, " ")
}
