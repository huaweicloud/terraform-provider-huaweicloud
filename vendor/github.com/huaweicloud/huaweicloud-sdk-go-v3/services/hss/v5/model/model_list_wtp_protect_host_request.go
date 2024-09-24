package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListWtpProtectHostRequest Request Object
type ListWtpProtectHostRequest struct {

	// Region Id
	Region string `json:"region"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 弹性公网IP
	PublicIp *string `json:"public_ip,omitempty"`

	// 私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 操作系统类别（linux，windows）   - linux : linux操作系统   - windows : windows操作系统
	OsType *string `json:"os_type,omitempty"`

	// 防护状态   - closed : 未开启   - opened : 防护中
	ProtectStatus *string `json:"protect_status,omitempty"`

	// 客户端状态   - not_installed : agent未安装   - online : agent在线   - offline : agent不在线
	AgentStatus *string `json:"agent_status,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListWtpProtectHostRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListWtpProtectHostRequest struct{}"
	}

	return strings.Join([]string{"ListWtpProtectHostRequest", string(data)}, " ")
}
