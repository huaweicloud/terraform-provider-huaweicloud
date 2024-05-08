package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProtectionPolicyInfo struct {

	// 策略ID
	PolicyId *string `json:"policy_id,omitempty"`

	// 策略名称
	PolicyName *string `json:"policy_name,omitempty"`

	// 防护动作，包含如下2种。   - alarm_and_isolation ：告警并自动隔离。   - alarm_only ：仅告警。
	ProtectionMode *string `json:"protection_mode,omitempty"`

	// 是否开启诱饵防护，包含如下1种, 默认为开启防护诱饵防护。   - opened ：开启。   - closed ：关闭。
	BaitProtectionStatus *string `json:"bait_protection_status,omitempty"`

	// 是否开启动态诱饵防护，包含如下2种, 默认为关闭动态诱饵防护。   - opened ：开启。   - closed ：关闭。
	DeployMode *string `json:"deploy_mode,omitempty"`

	// 防护目录
	ProtectionDirectory *string `json:"protection_directory,omitempty"`

	// 防护文件类型，例如：docx，txt，avi
	ProtectionType *string `json:"protection_type,omitempty"`

	// 排除目录，选填
	ExcludeDirectory *string `json:"exclude_directory,omitempty"`

	// 是否运行时检测，包含如下2种，暂时只有关闭一种状态，为保留字段。   - opened ：开启。   - closed ：关闭。
	RuntimeDetectionStatus *string `json:"runtime_detection_status,omitempty"`

	// 运行时检测目录，现在为保留字段
	RuntimeDetectionDirectory *string `json:"runtime_detection_directory,omitempty"`

	// 关联server个数
	CountAssociatedServer *int32 `json:"count_associated_server,omitempty"`

	// 操作系统类型。 - Linux - Windows
	OperatingSystem *string `json:"operating_system,omitempty"`

	// 进程白名单
	ProcessWhitelist *[]TrustProcessInfo `json:"process_whitelist,omitempty"`

	// 是否为默认策略，包含如下2种。   - 0 ：非默认策略。   - 1 ：默认策略
	DefaultPolicy *int32 `json:"default_policy,omitempty"`
}

func (o ProtectionPolicyInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionPolicyInfo struct{}"
	}

	return strings.Join([]string{"ProtectionPolicyInfo", string(data)}, " ")
}
