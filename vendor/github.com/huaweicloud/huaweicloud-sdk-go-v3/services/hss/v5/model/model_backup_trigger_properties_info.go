package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 调度器属性
type BackupTriggerPropertiesInfo struct {

	// 调度器的调度策略，长度限制为10240个字符，参照iCalendar RFC 2445规范，但仅支持FREQ、BYDAY、BYHOUR、BYMINUTE四个参数，其中FREQ仅支持WEEKLY和DAILY，BYDAY支持一周七天（MO、TU、WE、TH、FR、SA、SU），BYHOUR支持0-23小时，BYMINUTE支持0-59分钟，并且时间点间隔不能小于一小时，一个备份策略可以同时设置多个备份时间点，一天最多可以设置24个时间点。
	Pattern *[]string `json:"pattern,omitempty"`

	// 调度器开始时间，例如：2020-01-08 09:59:49
	StartTime *string `json:"start_time,omitempty"`
}

func (o BackupTriggerPropertiesInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupTriggerPropertiesInfo struct{}"
	}

	return strings.Join([]string{"BackupTriggerPropertiesInfo", string(data)}, " ")
}
