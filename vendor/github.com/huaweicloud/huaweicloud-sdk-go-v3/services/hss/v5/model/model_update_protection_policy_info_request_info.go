package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateProtectionPolicyInfoRequestInfo struct {

	// 策略ID
	PolicyId *string `json:"policy_id,omitempty"`

	// 策略名称
	PolicyName *string `json:"policy_name,omitempty"`

	// 防护动作，包含如下2种。   - alarm_and_isolation ：告警并自动隔离。   - alarm_only ：仅告警。
	ProtectionMode *string `json:"protection_mode,omitempty"`

	// 是否开启诱饵防护，包含如下1种, 默认为开启防护诱饵防护。   - opened ：开启。   - closed ：关闭。
	BaitProtectionStatus *string `json:"bait_protection_status,omitempty"`

	// 防护目录
	ProtectionDirectory *string `json:"protection_directory,omitempty"`

	// 防护文件类型
	ProtectionType *string `json:"protection_type,omitempty"`

	// 排除目录
	ExcludeDirectory *string `json:"exclude_directory,omitempty"`

	// 关联server
	AgentIdList *[]string `json:"agent_id_list,omitempty"`

	// 操作系统，包含如下：   - Windows : 无需处理   - Linux : 已忽略
	OperatingSystem *string `json:"operating_system,omitempty"`

	// 是否运行时检测，包含如下2种，暂时只有关闭一种状态，为保留字段。   - opened ：开启。   - closed ：关闭。
	RuntimeDetectionStatus *string `json:"runtime_detection_status,omitempty"`
}

func (o UpdateProtectionPolicyInfoRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProtectionPolicyInfoRequestInfo struct{}"
	}

	return strings.Join([]string{"UpdateProtectionPolicyInfoRequestInfo", string(data)}, " ")
}
