package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type TriggerCondition struct {

	// 指标查询模式。 - “AOM”：AOM原生 - “PROM”：AOM Prometheus - “NATIVE_PROM”：原生Prometheus
	MetricQueryMode TriggerConditionMetricQueryMode `json:"metric_query_mode"`

	// 指标命名空间。
	MetricNamespace string `json:"metric_namespace"`

	// 指标名称。
	MetricName string `json:"metric_name"`

	// 指标单位。
	MetricUnit string `json:"metric_unit"`

	// 指标维度。
	MetricLabels []string `json:"metric_labels"`

	// Prometheus语句。
	Promql string `json:"promql"`

	// Prometheus语句模板。
	PromqlExpr *[]string `json:"promql_expr,omitempty"`

	// 连续周期个数。
	TriggerTimes *string `json:"trigger_times,omitempty"`

	// 检查频率周期。 - 当trigger_type 为“HOURLY”时，填“” - 当trigger_type为“DAILY”时，格式为：“小时” 例如 每天凌晨三点\"03:00\" - 当trigger_type为“WEEKLY”时，格式为：“星期 小时”例如每周一凌晨三点 “1 03:00” - 当trigger_type为“CRON”时，格式为 标准CRON表达式 - 当trigger_type为“FIXED_RATE”时，秒的取值为15s，30s，分钟为 1~59，小时为 1~24。例如：“15s”，“30s”，“1min”，“1h”
	TriggerInterval *string `json:"trigger_interval,omitempty"`

	// 触发频率的类型： - “FIXED_RATE”：固定间隔 - “HOURLY”：每小时 - “DAILY”：每天 - “WEEKLY”：每周 - “CRON”：Cron表达式
	TriggerType *TriggerConditionTriggerType `json:"trigger_type,omitempty"`

	// Prometheus原生监控时长。
	PromqlFor *string `json:"promql_for,omitempty"`

	// 统计方式： - average - minimum - maximum - sum - sampleCount
	AggregationType *string `json:"aggregation_type,omitempty"`

	// 判断条件：“>”,“<”,“=”,“>=”,“<=”
	Operator *string `json:"operator,omitempty"`

	// 键值对形式，键为告警级别，值为告警阈值
	Thresholds map[string]string `json:"thresholds,omitempty"`

	// 统计周期。 - “15s” - “30s” - “1m” - “5m” - “15m” - “1h”
	AggregationWindow *string `json:"aggregation_window,omitempty"`

	Cmdb *CmdbInfo `json:"cmdb,omitempty"`

	// 查询筛选条件。
	QueryMatch *string `json:"query_match,omitempty"`

	// 查询参数
	QueryParam string `json:"query_param"`

	// 监控层级。
	AomMonitorLevel *string `json:"aom_monitor_level,omitempty"`

	// 聚合方式。 - “by”：不分组 - “avg” - “max” - “min” - “sum”
	AggregateType *TriggerConditionAggregateType `json:"aggregate_type,omitempty"`

	// 当配置方式为全量指标时可选择的指标运算方式。 - “single”：单个指标进行运算 - “mix”：多个指标进行混合运算
	MetricStatisticMethod *TriggerConditionMetricStatisticMethod `json:"metric_statistic_method,omitempty"`

	// 混合运算的表达式。
	Expression *string `json:"expression,omitempty"`

	// 混合运算的promQL。
	MixPromql *string `json:"mix_promql,omitempty"`
}

func (o TriggerCondition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TriggerCondition struct{}"
	}

	return strings.Join([]string{"TriggerCondition", string(data)}, " ")
}

type TriggerConditionMetricQueryMode struct {
	value string
}

type TriggerConditionMetricQueryModeEnum struct {
	AOM         TriggerConditionMetricQueryMode
	PROM        TriggerConditionMetricQueryMode
	NATIVE_PROM TriggerConditionMetricQueryMode
}

func GetTriggerConditionMetricQueryModeEnum() TriggerConditionMetricQueryModeEnum {
	return TriggerConditionMetricQueryModeEnum{
		AOM: TriggerConditionMetricQueryMode{
			value: "AOM",
		},
		PROM: TriggerConditionMetricQueryMode{
			value: "PROM",
		},
		NATIVE_PROM: TriggerConditionMetricQueryMode{
			value: "NATIVE_PROM",
		},
	}
}

func (c TriggerConditionMetricQueryMode) Value() string {
	return c.value
}

func (c TriggerConditionMetricQueryMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TriggerConditionMetricQueryMode) UnmarshalJSON(b []byte) error {
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

type TriggerConditionTriggerType struct {
	value string
}

type TriggerConditionTriggerTypeEnum struct {
	FIXED_RATE TriggerConditionTriggerType
	HOURLY     TriggerConditionTriggerType
	DAILY      TriggerConditionTriggerType
	WEEKLY     TriggerConditionTriggerType
	CRON       TriggerConditionTriggerType
}

func GetTriggerConditionTriggerTypeEnum() TriggerConditionTriggerTypeEnum {
	return TriggerConditionTriggerTypeEnum{
		FIXED_RATE: TriggerConditionTriggerType{
			value: "FIXED_RATE",
		},
		HOURLY: TriggerConditionTriggerType{
			value: "HOURLY",
		},
		DAILY: TriggerConditionTriggerType{
			value: "DAILY",
		},
		WEEKLY: TriggerConditionTriggerType{
			value: "WEEKLY",
		},
		CRON: TriggerConditionTriggerType{
			value: "CRON",
		},
	}
}

func (c TriggerConditionTriggerType) Value() string {
	return c.value
}

func (c TriggerConditionTriggerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TriggerConditionTriggerType) UnmarshalJSON(b []byte) error {
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

type TriggerConditionAggregateType struct {
	value string
}

type TriggerConditionAggregateTypeEnum struct {
	BY  TriggerConditionAggregateType
	AVG TriggerConditionAggregateType
	MAX TriggerConditionAggregateType
	MIN TriggerConditionAggregateType
	SUM TriggerConditionAggregateType
}

func GetTriggerConditionAggregateTypeEnum() TriggerConditionAggregateTypeEnum {
	return TriggerConditionAggregateTypeEnum{
		BY: TriggerConditionAggregateType{
			value: "by",
		},
		AVG: TriggerConditionAggregateType{
			value: "avg",
		},
		MAX: TriggerConditionAggregateType{
			value: "max",
		},
		MIN: TriggerConditionAggregateType{
			value: "min",
		},
		SUM: TriggerConditionAggregateType{
			value: "sum",
		},
	}
}

func (c TriggerConditionAggregateType) Value() string {
	return c.value
}

func (c TriggerConditionAggregateType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TriggerConditionAggregateType) UnmarshalJSON(b []byte) error {
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

type TriggerConditionMetricStatisticMethod struct {
	value string
}

type TriggerConditionMetricStatisticMethodEnum struct {
	SINGLE TriggerConditionMetricStatisticMethod
	MIX    TriggerConditionMetricStatisticMethod
}

func GetTriggerConditionMetricStatisticMethodEnum() TriggerConditionMetricStatisticMethodEnum {
	return TriggerConditionMetricStatisticMethodEnum{
		SINGLE: TriggerConditionMetricStatisticMethod{
			value: "single",
		},
		MIX: TriggerConditionMetricStatisticMethod{
			value: "mix",
		},
	}
}

func (c TriggerConditionMetricStatisticMethod) Value() string {
	return c.value
}

func (c TriggerConditionMetricStatisticMethod) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TriggerConditionMetricStatisticMethod) UnmarshalJSON(b []byte) error {
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
