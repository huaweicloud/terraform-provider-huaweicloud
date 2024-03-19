package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ContainerNodeInfo 容器节点列表
type ContainerNodeInfo struct {

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// 节点名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器状态，包含如下4种。   - ACTIVE ：正在运行。   - SHUTOFF ：关机。   - BUILDING ：创建中。   - ERROR ：故障。
	HostStatus *string `json:"host_status,omitempty"`

	// Agent状态，包含如下3种。   - not_installed ：未安装。   - online ：在线。   - offline ：离线。
	AgentStatus *string `json:"agent_status,omitempty"`

	// 防护状态，包含如下2种。   - closed ：关闭。   - opened ：开启。
	ProtectStatus *string `json:"protect_status,omitempty"`

	// 防护是否中断
	ProtectInterrupt *bool `json:"protect_interrupt,omitempty"`

	// 标签：用来识别cce容器节点和自建  - cce：cce节点  - self：自建节点  - other：其他节点
	ContainerTags *string `json:"container_tags,omitempty"`

	// 私有IP地址
	PrivateIp *string `json:"private_ip,omitempty"`

	// 弹性公网IP地址
	PublicIp *string `json:"public_ip,omitempty"`

	// 主机安全配额ID（UUID）
	ResourceId *string `json:"resource_id,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 所属企业项目名称
	EnterpriseProjectName *string `json:"enterprise_project_name,omitempty"`

	// 云主机安全检测结果，包含如下4种。 - undetected ：未检测。 - clean ：无风险。 - risk ：有风险。 - scanning ：检测中。
	DetectResult *string `json:"detect_result,omitempty"`

	// 资产风险
	Asset *int32 `json:"asset,omitempty"`

	// 漏洞风险
	Vulnerability *int32 `json:"vulnerability,omitempty"`

	// 入侵风险
	Intrusion *int32 `json:"intrusion,omitempty"`

	// 策略组ID
	PolicyGroupId *string `json:"policy_group_id,omitempty"`

	// 策略组名称
	PolicyGroupName *string `json:"policy_group_name,omitempty"`
}

func (o ContainerNodeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContainerNodeInfo struct{}"
	}

	return strings.Join([]string{"ContainerNodeInfo", string(data)}, " ")
}
