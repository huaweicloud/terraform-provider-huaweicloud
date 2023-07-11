package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ProtectionProxyInfoRequestInfo 创建防护策略。若新建防护策略，则protection_policy_id为空，create_protection_policy必选
type ProtectionProxyInfoRequestInfo struct {

	// 策略ID，新建策略可不填
	PolicyId *string `json:"policy_id,omitempty"`

	// 策略名称，新建防护策略则必填
	PolicyName *string `json:"policy_name,omitempty"`

	// 防护动作，新建防护策略则必填。包含如下：   - alarm_and_isolation ：告警并自动隔离。   - alarm_only ：仅告警。
	ProtectionMode *string `json:"protection_mode,omitempty"`

	// 是否开启诱饵防护，新建防护策略则必填。包含如下1种, 默认为开启防护诱饵防护。   - opened ：开启。   - closed ：关闭。
	BaitProtectionStatus *string `json:"bait_protection_status,omitempty"`

	// 防护目录，新建防护策略则必填
	ProtectionDirectory *string `json:"protection_directory,omitempty"`

	// 防护类型，新建防护策略则必填
	ProtectionType *string `json:"protection_type,omitempty"`

	// 排除目录，可选填
	ExcludeDirectory *string `json:"exclude_directory,omitempty"`

	// 是否运行时检测，选填。包含如下2种，暂时只有关闭一种状态，为保留字段。   - opened ：开启。   - closed ：关闭。
	RuntimeDetectionStatus *string `json:"runtime_detection_status,omitempty"`

	// 操作系统，新建防护策略则必填。包含如下：   - Windows : Windows系统   - Linux : Linux系统
	OperatingSystem *string `json:"operating_system,omitempty"`

	// 进程白名单
	ProcessWhitelist *[]TrustProcessInfo `json:"process_whitelist,omitempty"`
}

func (o ProtectionProxyInfoRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionProxyInfoRequestInfo struct{}"
	}

	return strings.Join([]string{"ProtectionProxyInfoRequestInfo", string(data)}, " ")
}
