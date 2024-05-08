package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeEventRequestInfo 操作事件详情
type ChangeEventRequestInfo struct {

	// 处理方式，包含如下:   - mark_as_handled : 手动处理   - ignore : 忽略   - add_to_alarm_whitelist : 加入告警白名单   - add_to_login_whitelist : 加入登录白名单   - isolate_and_kill : 隔离查杀   - unhandle : 取消手动处理   - do_not_ignore : 取消忽略   - remove_from_alarm_whitelist : 删除告警白名单   - remove_from_login_whitelist : 删除登录白名单   - do_not_isolate_or_kill : 取消隔离查杀
	OperateType string `json:"operate_type"`

	// 备注信息，已处理的告警才有
	Handler *string `json:"handler,omitempty"`

	// 操作的事件列表
	OperateEventList []OperateEventRequestInfo `json:"operate_event_list"`

	// 用户自定义告警白名单规则列表
	EventWhiteRuleList *[]EventWhiteRuleListRequestInfo `json:"event_white_rule_list,omitempty"`
}

func (o ChangeEventRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeEventRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeEventRequestInfo", string(data)}, " ")
}
