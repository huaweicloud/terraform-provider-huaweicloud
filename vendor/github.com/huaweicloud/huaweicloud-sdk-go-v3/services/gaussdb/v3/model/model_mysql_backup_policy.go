package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlBackupPolicy struct {
	// 备份时间段。自动备份将在该时间段内触发。取值范围：非空，格式必须为hh:mm-HH:MM且有效，当前时间指UTC时间。HH取值必须比hh大1。mm和MM取值必须相同，且取值必须为00。取值示例：21:00-22:00

	StartTime string `json:"start_time"`
	// 备份文件的保留天数。

	KeepDays int32 `json:"keep_days"`
	// 备份周期配置。自动备份将在每星期指定的天进行。取值范围：格式为逗号隔开的数字，数字代表星期。取值示例：1,2,3,4则表示备份周期配置为星期一、星期二、星期三和星期四。

	Period string `json:"period"`
	// 1级备份保留数量，默认值为0。当一级备份开关开启时，该参数值有效。取值：0或1

	RetentionNumBackupLevel1 *int32 `json:"retention_num_backup_level1,omitempty"`
}

func (o MysqlBackupPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlBackupPolicy struct{}"
	}

	return strings.Join([]string{"MysqlBackupPolicy", string(data)}, " ")
}
