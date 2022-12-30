package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListHostRaspProtectHistoryInfoRequest struct {

	// Region Id
	Region string `json:"region"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// Host Id
	HostId string `json:"host_id"`

	// 起始时间
	StartTime int64 `json:"start_time"`

	// 终止时间
	EndTime int64 `json:"end_time"`

	// limit
	Limit int32 `json:"limit"`

	// offset
	Offset int32 `json:"offset"`

	// 告警级别
	AlarmLevel *int32 `json:"alarm_level,omitempty"`

	// 威胁等级   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危   - Critical : 危急
	Severity *string `json:"severity,omitempty"`

	// 防护状态   - closed : 未开启   - opened : 防护中
	ProtectStatus *string `json:"protect_status,omitempty"`
}

func (o ListHostRaspProtectHistoryInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostRaspProtectHistoryInfoRequest struct{}"
	}

	return strings.Join([]string{"ListHostRaspProtectHistoryInfoRequest", string(data)}, " ")
}
