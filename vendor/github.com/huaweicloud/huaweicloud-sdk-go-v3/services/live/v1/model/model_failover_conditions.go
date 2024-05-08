package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// FailoverConditions 主备切换配置
type FailoverConditions struct {

	// 入流停止的时长阈值。到达此阈值后，自动触发主备切换  单位：毫秒，取值范围：0 - 3600000
	InputLossThresholdMsec *int32 `json:"input_loss_threshold_msec,omitempty"`

	// 以主入流URL为第一优先级（PRIMARY）或主备URL平等切换（EQUAL）  如果为平等切换时使用的是备URL，无需手工切换到主URL
	InputPreference *FailoverConditionsInputPreference `json:"input_preference,omitempty"`
}

func (o FailoverConditions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FailoverConditions struct{}"
	}

	return strings.Join([]string{"FailoverConditions", string(data)}, " ")
}

type FailoverConditionsInputPreference struct {
	value string
}

type FailoverConditionsInputPreferenceEnum struct {
	EQUAL   FailoverConditionsInputPreference
	PRIMARY FailoverConditionsInputPreference
}

func GetFailoverConditionsInputPreferenceEnum() FailoverConditionsInputPreferenceEnum {
	return FailoverConditionsInputPreferenceEnum{
		EQUAL: FailoverConditionsInputPreference{
			value: "EQUAL",
		},
		PRIMARY: FailoverConditionsInputPreference{
			value: "PRIMARY",
		},
	}
}

func (c FailoverConditionsInputPreference) Value() string {
	return c.value
}

func (c FailoverConditionsInputPreference) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *FailoverConditionsInputPreference) UnmarshalJSON(b []byte) error {
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
