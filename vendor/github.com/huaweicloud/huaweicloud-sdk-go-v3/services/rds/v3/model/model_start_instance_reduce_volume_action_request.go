package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// StartInstanceReduceVolumeActionRequest Request Object
type StartInstanceReduceVolumeActionRequest struct {

	// 语言
	XLanguage *StartInstanceReduceVolumeActionRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ReduceVolumeRequestBody `json:"body,omitempty"`
}

func (o StartInstanceReduceVolumeActionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceReduceVolumeActionRequest struct{}"
	}

	return strings.Join([]string{"StartInstanceReduceVolumeActionRequest", string(data)}, " ")
}

type StartInstanceReduceVolumeActionRequestXLanguage struct {
	value string
}

type StartInstanceReduceVolumeActionRequestXLanguageEnum struct {
	ZH_CN StartInstanceReduceVolumeActionRequestXLanguage
	EN_US StartInstanceReduceVolumeActionRequestXLanguage
}

func GetStartInstanceReduceVolumeActionRequestXLanguageEnum() StartInstanceReduceVolumeActionRequestXLanguageEnum {
	return StartInstanceReduceVolumeActionRequestXLanguageEnum{
		ZH_CN: StartInstanceReduceVolumeActionRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: StartInstanceReduceVolumeActionRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c StartInstanceReduceVolumeActionRequestXLanguage) Value() string {
	return c.value
}

func (c StartInstanceReduceVolumeActionRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *StartInstanceReduceVolumeActionRequestXLanguage) UnmarshalJSON(b []byte) error {
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
