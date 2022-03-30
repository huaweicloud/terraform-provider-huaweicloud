package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 秒级监控修改请求
type TaurusModifyInstanceMonitorRequestBody struct {
	// 实例秒级监控开关。为true时表示开启，为false时表示关闭。

	MonitorSwitch bool `json:"monitor_switch"`
	// 采集周期，仅在monitor_switch为true时生效。默认为5s。monitor_switch为false时，不传该参数。  取值： 1：采集周期为1s； 5：采集周期为5s。

	Period *TaurusModifyInstanceMonitorRequestBodyPeriod `json:"period,omitempty"`
}

func (o TaurusModifyInstanceMonitorRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaurusModifyInstanceMonitorRequestBody struct{}"
	}

	return strings.Join([]string{"TaurusModifyInstanceMonitorRequestBody", string(data)}, " ")
}

type TaurusModifyInstanceMonitorRequestBodyPeriod struct {
	value int32
}

type TaurusModifyInstanceMonitorRequestBodyPeriodEnum struct {
	E_1 TaurusModifyInstanceMonitorRequestBodyPeriod
	E_5 TaurusModifyInstanceMonitorRequestBodyPeriod
}

func GetTaurusModifyInstanceMonitorRequestBodyPeriodEnum() TaurusModifyInstanceMonitorRequestBodyPeriodEnum {
	return TaurusModifyInstanceMonitorRequestBodyPeriodEnum{
		E_1: TaurusModifyInstanceMonitorRequestBodyPeriod{
			value: 1,
		}, E_5: TaurusModifyInstanceMonitorRequestBodyPeriod{
			value: 5,
		},
	}
}

func (c TaurusModifyInstanceMonitorRequestBodyPeriod) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaurusModifyInstanceMonitorRequestBodyPeriod) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(int32)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to int32 error")
	}
}
