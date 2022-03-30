package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 自动备份策略
type MysqlBackupStrategy struct {
	// 自动备份开始时间段。自动备份将在该时间一个小时内触发。  取值范围：非空，格式必须为hh:mm-HH:MM且有效，当前时间指UTC时间。  1. HH取值必须比hh大1。 2. mm和MM取值必须相同，且取值必须为00。

	StartTime string `json:"start_time"`
	// 自动备份保留天数，取值范围：1-732

	KeepDays *string `json:"keep_days,omitempty"`
}

func (o MysqlBackupStrategy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlBackupStrategy struct{}"
	}

	return strings.Join([]string{"MysqlBackupStrategy", string(data)}, " ")
}
