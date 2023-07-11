package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateInstanceByEngineRequest Request Object
type CreateInstanceByEngineRequest struct {

	// 消息引擎。
	Engine CreateInstanceByEngineRequestEngine `json:"engine"`

	Body *CreateInstanceByEngineReq `json:"body,omitempty"`
}

func (o CreateInstanceByEngineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceByEngineRequest struct{}"
	}

	return strings.Join([]string{"CreateInstanceByEngineRequest", string(data)}, " ")
}

type CreateInstanceByEngineRequestEngine struct {
	value string
}

type CreateInstanceByEngineRequestEngineEnum struct {
	KAFKA CreateInstanceByEngineRequestEngine
}

func GetCreateInstanceByEngineRequestEngineEnum() CreateInstanceByEngineRequestEngineEnum {
	return CreateInstanceByEngineRequestEngineEnum{
		KAFKA: CreateInstanceByEngineRequestEngine{
			value: "kafka",
		},
	}
}

func (c CreateInstanceByEngineRequestEngine) Value() string {
	return c.value
}

func (c CreateInstanceByEngineRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateInstanceByEngineRequestEngine) UnmarshalJSON(b []byte) error {
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
