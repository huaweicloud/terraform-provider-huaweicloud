package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type NoDataCondition struct {

	// 无数据周期的个数。
	NoDataTimeframe *int32 `json:"no_data_timeframe,omitempty"`

	// 数据不足时，阈值规则的状态。 - “no_data”：数据不足并发送通知 - “alerting”：告警 - “ok”：正常 - “pre_state”：保持上一个状态
	NoDataAlertState *NoDataConditionNoDataAlertState `json:"no_data_alert_state,omitempty"`

	// 数据不足是否通知。
	NotifyNoData *bool `json:"notify_no_data,omitempty"`
}

func (o NoDataCondition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NoDataCondition struct{}"
	}

	return strings.Join([]string{"NoDataCondition", string(data)}, " ")
}

type NoDataConditionNoDataAlertState struct {
	value string
}

type NoDataConditionNoDataAlertStateEnum struct {
	NO_DATA   NoDataConditionNoDataAlertState
	ALERTING  NoDataConditionNoDataAlertState
	OK        NoDataConditionNoDataAlertState
	PRE_STATE NoDataConditionNoDataAlertState
}

func GetNoDataConditionNoDataAlertStateEnum() NoDataConditionNoDataAlertStateEnum {
	return NoDataConditionNoDataAlertStateEnum{
		NO_DATA: NoDataConditionNoDataAlertState{
			value: "no_data",
		},
		ALERTING: NoDataConditionNoDataAlertState{
			value: "alerting",
		},
		OK: NoDataConditionNoDataAlertState{
			value: "ok",
		},
		PRE_STATE: NoDataConditionNoDataAlertState{
			value: "pre_state",
		},
	}
}

func (c NoDataConditionNoDataAlertState) Value() string {
	return c.value
}

func (c NoDataConditionNoDataAlertState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NoDataConditionNoDataAlertState) UnmarshalJSON(b []byte) error {
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
