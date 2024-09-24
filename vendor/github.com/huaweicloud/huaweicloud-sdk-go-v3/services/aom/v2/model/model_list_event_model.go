package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListEventModel 事件或者告警元数据。
type ListEventModel struct {

	// 事件或者告警产生的时间，CST毫秒级时间戳。
	StartsAt *int64 `json:"starts_at,omitempty"`

	// 事件或者告警清除的时间，CST毫秒级时间戳，为0时表示未删除。
	EndsAt *int64 `json:"ends_at,omitempty"`

	// 告警自动清除时间。毫秒数，例如一分钟则填写为60000。默认清除时间为3天,对应数字为 4320 * 1000（即：3天 * 24小时 * 60分钟 * 1000毫秒）。
	Timeout *int64 `json:"timeout,omitempty"`

	// 事件或者告警的详细信息，为键值对形式。必须字段为：  - event_name：事件或者告警名称,类型为String；  - event_severity：事件级别枚举值。类型为String，四种类型 \"Critical\", \"Major\", \"Minor\", \"Info\"；  - event_type：事件类别枚举值。类型为String，event为告警事件，alarm为普通告警；  - resource_provider：事件对应云服务名称。类型为String；  - resource_type：事件对应资源类型。类型为String；  - resource_id：事件对应资源信息。类型为String。
	Metadata map[string]string `json:"metadata,omitempty"`

	// 事件或者告警附加字段，可以为空。
	Annotations map[string]interface{} `json:"annotations,omitempty"`

	// 事件或者告警预留字段，为空。
	AttachRule map[string]interface{} `json:"attach_rule,omitempty"`

	// 事件或者告警id，系统会自动生成，上报无须填写该字段。
	Id *string `json:"id,omitempty"`

	// 告警流水号。
	EventSn *string `json:"event_sn,omitempty"`

	// 事件到达系统时间，CST毫秒级时间戳。
	ArrivesAt *int64 `json:"arrives_at,omitempty"`

	// 事件或告警所属企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 开放告警策略
	Policy map[string]interface{} `json:"policy,omitempty"`
}

func (o ListEventModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEventModel struct{}"
	}

	return strings.Join([]string{"ListEventModel", string(data)}, " ")
}
