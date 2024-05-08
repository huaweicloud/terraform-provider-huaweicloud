package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type WtpProtectHostResponseInfo struct {

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 弹性公网IP
	PublicIp *string `json:"public_ip,omitempty"`

	// 私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 私有IPv6地址
	Ipv6 *string `json:"ipv6,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 操作系统位数
	OsBit *string `json:"os_bit,omitempty"`

	// 操作系统（linux，windows）
	OsType *string `json:"os_type,omitempty"`

	// 防护状态   - closed : 未开启   - opened : 防护中
	ProtectStatus *string `json:"protect_status,omitempty"`

	// 动态网页防篡改状态   - closed : 未开启   - opened : 防护中
	RaspProtectStatus *string `json:"rasp_protect_status,omitempty"`

	// 已防御篡改攻击次数
	AntiTamperingTimes *int64 `json:"anti_tampering_times,omitempty"`

	// 已发现篡改攻击
	DetectTamperingTimes *int64 `json:"detect_tampering_times,omitempty"`

	// 最近检测时间(ms)
	LastDetectTime *int64 `json:"last_detect_time,omitempty"`

	// 定时关闭防护开关状态   - opened : 开启   - closed : 未开启
	ScheduledShutdownStatus *string `json:"scheduled_shutdown_status,omitempty"`

	// Agent状态   - not_installed : agent未安装   - online : agent在线   - offline : agent不在线
	AgentStatus *string `json:"agent_status,omitempty"`
}

func (o WtpProtectHostResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WtpProtectHostResponseInfo struct{}"
	}

	return strings.Join([]string{"WtpProtectHostResponseInfo", string(data)}, " ")
}
