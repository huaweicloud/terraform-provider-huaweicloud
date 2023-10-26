package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ResizeEngineInstanceRequest Request Object
type ResizeEngineInstanceRequest struct {

	// 消息引擎。
	Engine ResizeEngineInstanceRequestEngine `json:"engine"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ResizeEngineInstanceReq `json:"body,omitempty"`
}

func (o ResizeEngineInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeEngineInstanceRequest struct{}"
	}

	return strings.Join([]string{"ResizeEngineInstanceRequest", string(data)}, " ")
}

type ResizeEngineInstanceRequestEngine struct {
	value string
}

type ResizeEngineInstanceRequestEngineEnum struct {
	KAFKA ResizeEngineInstanceRequestEngine
}

func GetResizeEngineInstanceRequestEngineEnum() ResizeEngineInstanceRequestEngineEnum {
	return ResizeEngineInstanceRequestEngineEnum{
		KAFKA: ResizeEngineInstanceRequestEngine{
			value: "kafka",
		},
	}
}

func (c ResizeEngineInstanceRequestEngine) Value() string {
	return c.value
}

func (c ResizeEngineInstanceRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResizeEngineInstanceRequestEngine) UnmarshalJSON(b []byte) error {
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
