package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type PreheatingResult struct {

	// 媒资URL。
	Url *string `json:"url,omitempty"`

	// 预热任务状态。  取值如下： - processing：处理中 - succeed：预热完成 - failed：预热失败
	Status *PreheatingResultStatus `json:"status,omitempty"`
}

func (o PreheatingResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PreheatingResult struct{}"
	}

	return strings.Join([]string{"PreheatingResult", string(data)}, " ")
}

type PreheatingResultStatus struct {
	value string
}

type PreheatingResultStatusEnum struct {
	PROCESSING PreheatingResultStatus
	SUCCEED    PreheatingResultStatus
	FAILED     PreheatingResultStatus
}

func GetPreheatingResultStatusEnum() PreheatingResultStatusEnum {
	return PreheatingResultStatusEnum{
		PROCESSING: PreheatingResultStatus{
			value: "PROCESSING",
		},
		SUCCEED: PreheatingResultStatus{
			value: "SUCCEED",
		},
		FAILED: PreheatingResultStatus{
			value: "FAILED",
		},
	}
}

func (c PreheatingResultStatus) Value() string {
	return c.value
}

func (c PreheatingResultStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PreheatingResultStatus) UnmarshalJSON(b []byte) error {
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
