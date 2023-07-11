package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OperationDefinitionRequestInfo 调度参数
type OperationDefinitionRequestInfo struct {

	// 保留日备个数，该备份不受保留最大备份数限制。取值为0到100。若选择该参数，则timezone 也必选。最小值：0,最大值：100
	DayBackups *int32 `json:"day_backups,omitempty"`

	// 单个备份对象自动备份的最大备份数。取值为-1或0-99999。-1代表不按备份数清理。若该字段和retention_duration_days字段同时为空，备份会永久保留。最小值：1,最大值：99999,缺省值：-1
	MaxBackups *int32 `json:"max_backups,omitempty"`

	// 保留月备个数，该备份不受保留最大备份数限制。取值为0到100。若选择该参数，则timezone 也必选。最小值：0, 最大值：100
	MonthBackups *int32 `json:"month_backups,omitempty"`

	// 备份保留时长，单位天。最长支持99999天。-1代表不按时间清理。若该字段和max_backups 参数同时为空，备份会永久保留。最小值：1, 最大值：99999, 缺省值：-1
	RetentionDurationDays *int32 `json:"retention_duration_days,omitempty"`

	// 用户所在时区,格式形如UTC+08:00,若没有选择年备，月备，周备，日备中任一参数，则不能选择该参数。
	Timezone *string `json:"timezone,omitempty"`

	// 保留周备个数，该备份不受保留最大备份数限制。取值为0到100。若选择该参数，则timezone 也必选。
	WeekBackups *int32 `json:"week_backups,omitempty"`

	// 保留年备个数，该备份不受保留最大备份数限制。取值为0到100。若选择该参数，则timezone 也必选。最小值：0，最大值：100
	YearBackups *int32 `json:"year_backups,omitempty"`
}

func (o OperationDefinitionRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OperationDefinitionRequestInfo struct{}"
	}

	return strings.Join([]string{"OperationDefinitionRequestInfo", string(data)}, " ")
}
