package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 事件或者告警元数据。
type EventModel struct {

	// 事件或者告警产生的时间，CST毫秒级时间戳。
	StartsAt *int64 `json:"starts_at,omitempty"`

	// 事件或者告警清除的时间，CST毫秒级时间戳，为0时表示未删除。
	EndsAt *int64 `json:"ends_at,omitempty"`

	// 告警自动清除时间。毫秒数，例如一分钟则填写为60000。默认清除时间为3天,对应数字为 4320 * 1000（即：3天 * 24小时 * 60分钟 * 1000毫秒）。
	Timeout *int64 `json:"timeout,omitempty"`

	// 事件或者告警的详细信息，为键值对形式。必须字段为： - event_name：事件或者告警名称,类型为String； - event_severity：事件级别枚举值。类型为String，四种类型 \"Critical\", \"Major\", \"Minor\", \"Info\"； - event_type：事件类别枚举值。类型为String，event为普通告警，alarm为告警事件； - resource_provider：事件对应云服务名称。类型为String； - resource_type：事件对应资源类型。类型为String； - resource_id：事件对应资源信息。类型为String。
	Metadata *interface{} `json:"metadata,omitempty"`

	// 事件或者告警附加字段，可以为空。
	Annotations *interface{} `json:"annotations,omitempty"`

	// 事件或者告警预留字段，为空。
	AttachRule *interface{} `json:"attach_rule,omitempty"`

	// 事件或者告警id，系统会自动生成，上报无须填写该字段。
	Id *string `json:"id,omitempty"`
}

func (o EventModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventModel struct{}"
	}

	return strings.Join([]string{"EventModel", string(data)}, " ")
}
