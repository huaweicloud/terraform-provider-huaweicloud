package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 备份策略信息。
type BackupPolicy struct {
	// 指定已生成的备份文件可以保存的天数。取值范围：1～732。

	KeepDays int32 `json:"keep_days"`
	// 备份时间段。自动备份将在该时间段内触发。 取值范围：格式必须为hh:mm-HH:MM且有效，当前时间指UTC时间。

	StartTime *string `json:"start_time,omitempty"`
	// 备份周期配置。自动备份将在每星期指定的天进行。 取值范围：格式为逗号隔开的数字，数字代表星期。

	Period *string `json:"period,omitempty"`
	// 1级备份保留数量。当一级备份开关开启时，返回此参数。

	RetentionNumBackupLevel1 *int32 `json:"retention_num_backup_level1,omitempty"`
}

func (o BackupPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupPolicy struct{}"
	}

	return strings.Join([]string{"BackupPolicy", string(data)}, " ")
}
