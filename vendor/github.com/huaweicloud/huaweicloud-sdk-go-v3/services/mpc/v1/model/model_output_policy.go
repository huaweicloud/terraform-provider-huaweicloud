package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type OutputPolicy struct {

	// 输出策略。  取值如下： - discard - transcode  >- 当视频参数中的“output_policy”为\"discard\"，且音频参数中的“output_policy”为“transcode”时，表示只输出音频。 >- 当视频参数中的“output_policy”为\"transcode\"，且音频参数中的“output_policy”为“discard”时，表示只输出视频。 >- 同时为\"discard\"时不合法。 >- 同时为“transcode”时，表示输出音视频。
	OutputPolicy *OutputPolicyOutputPolicy `json:"output_policy,omitempty"`
}

func (o OutputPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OutputPolicy struct{}"
	}

	return strings.Join([]string{"OutputPolicy", string(data)}, " ")
}

type OutputPolicyOutputPolicy struct {
	value string
}

type OutputPolicyOutputPolicyEnum struct {
	TRANSCODE OutputPolicyOutputPolicy
	DISCARD   OutputPolicyOutputPolicy
	COPY      OutputPolicyOutputPolicy
}

func GetOutputPolicyOutputPolicyEnum() OutputPolicyOutputPolicyEnum {
	return OutputPolicyOutputPolicyEnum{
		TRANSCODE: OutputPolicyOutputPolicy{
			value: "transcode",
		},
		DISCARD: OutputPolicyOutputPolicy{
			value: "discard",
		},
		COPY: OutputPolicyOutputPolicy{
			value: "copy",
		},
	}
}

func (c OutputPolicyOutputPolicy) Value() string {
	return c.value
}

func (c OutputPolicyOutputPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OutputPolicyOutputPolicy) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
