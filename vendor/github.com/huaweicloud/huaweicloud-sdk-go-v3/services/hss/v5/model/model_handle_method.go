package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HandleMethod 处理方式，已处理的告警才有，包含如下:   - mark_as_handled : 手动处理   - ignore : 忽略   - add_to_alarm_whitelist : 加入告警白名单   - add_to_login_whitelist : 加入登录白名单   - isolate_and_kill : 隔离查杀
type HandleMethod struct {
}

func (o HandleMethod) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HandleMethod struct{}"
	}

	return strings.Join([]string{"HandleMethod", string(data)}, " ")
}
