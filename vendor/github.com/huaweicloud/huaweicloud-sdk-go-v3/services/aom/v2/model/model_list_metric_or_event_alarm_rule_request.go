package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListMetricOrEventAlarmRuleRequest Request Object
type ListMetricOrEventAlarmRuleRequest struct {

	// 告警规则名称。
	Name *string `json:"name,omitempty"`

	// 用于限制本次返回的结果数据条数。
	Limit *string `json:"limit,omitempty"`

	// 分页查询起始位置，为非负整数。
	Offset *string `json:"offset,omitempty"`

	// 根据告警规则名称或者告警创建时间排序。 - alarm_rule_name.asc - alarm_create_time.desc
	SortBy *string `json:"sort_by,omitempty"`

	// 事件告警规则事件来源。 - “RDS” - “EVS” - “CCE” - “LTS” - “AOM”
	EventSource *string `json:"event_source,omitempty"`

	// 事件告警级别。 - “Critical\" - “Major” - “Minor” - “Info”
	EventSeverity *string `json:"event_severity,omitempty"`

	// 告警规则状态。 - “OK”：正常 - “alarm”：超限阈值 - “Effective”：生效中 - “Invalid”：停用中
	AlarmRuleStatus *ListMetricOrEventAlarmRuleRequestAlarmRuleStatus `json:"alarm_rule_status,omitempty"`

	// 告警规则类型。 - “metric”：指标告警规则 - “event”：事件告警规则
	AlarmRuleType *ListMetricOrEventAlarmRuleRequestAlarmRuleType `json:"alarm_rule_type,omitempty"`

	// Prometheus实例id。
	PromInstanceId *string `json:"prom_instance_id,omitempty"`

	// 绑定的告警行动规则名称。
	BindNotificationRuleId *string `json:"bind_notification_rule_id,omitempty"`

	// CCE集群id。
	RelatedCceClusters *string `json:"related_cce_clusters,omitempty"`

	// 企业项目id。  - 查询单个企业项目下实例，填写企业项目id。  - 查询所有企业项目下实例，填写“all_granted_eps”。
	EnterpriseProjectId *string `json:"Enterprise-Project-Id,omitempty"`
}

func (o ListMetricOrEventAlarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMetricOrEventAlarmRuleRequest struct{}"
	}

	return strings.Join([]string{"ListMetricOrEventAlarmRuleRequest", string(data)}, " ")
}

type ListMetricOrEventAlarmRuleRequestAlarmRuleStatus struct {
	value string
}

type ListMetricOrEventAlarmRuleRequestAlarmRuleStatusEnum struct {
	OK        ListMetricOrEventAlarmRuleRequestAlarmRuleStatus
	ALARM     ListMetricOrEventAlarmRuleRequestAlarmRuleStatus
	EFFECTIVE ListMetricOrEventAlarmRuleRequestAlarmRuleStatus
	INVALID   ListMetricOrEventAlarmRuleRequestAlarmRuleStatus
}

func GetListMetricOrEventAlarmRuleRequestAlarmRuleStatusEnum() ListMetricOrEventAlarmRuleRequestAlarmRuleStatusEnum {
	return ListMetricOrEventAlarmRuleRequestAlarmRuleStatusEnum{
		OK: ListMetricOrEventAlarmRuleRequestAlarmRuleStatus{
			value: "OK",
		},
		ALARM: ListMetricOrEventAlarmRuleRequestAlarmRuleStatus{
			value: "alarm",
		},
		EFFECTIVE: ListMetricOrEventAlarmRuleRequestAlarmRuleStatus{
			value: "Effective",
		},
		INVALID: ListMetricOrEventAlarmRuleRequestAlarmRuleStatus{
			value: "Invalid",
		},
	}
}

func (c ListMetricOrEventAlarmRuleRequestAlarmRuleStatus) Value() string {
	return c.value
}

func (c ListMetricOrEventAlarmRuleRequestAlarmRuleStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListMetricOrEventAlarmRuleRequestAlarmRuleStatus) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type ListMetricOrEventAlarmRuleRequestAlarmRuleType struct {
	value string
}

type ListMetricOrEventAlarmRuleRequestAlarmRuleTypeEnum struct {
	METRIC ListMetricOrEventAlarmRuleRequestAlarmRuleType
	EVENT  ListMetricOrEventAlarmRuleRequestAlarmRuleType
}

func GetListMetricOrEventAlarmRuleRequestAlarmRuleTypeEnum() ListMetricOrEventAlarmRuleRequestAlarmRuleTypeEnum {
	return ListMetricOrEventAlarmRuleRequestAlarmRuleTypeEnum{
		METRIC: ListMetricOrEventAlarmRuleRequestAlarmRuleType{
			value: "metric",
		},
		EVENT: ListMetricOrEventAlarmRuleRequestAlarmRuleType{
			value: "event",
		},
	}
}

func (c ListMetricOrEventAlarmRuleRequestAlarmRuleType) Value() string {
	return c.value
}

func (c ListMetricOrEventAlarmRuleRequestAlarmRuleType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListMetricOrEventAlarmRuleRequestAlarmRuleType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
