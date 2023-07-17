package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowEngineInstanceExtendProductInfoRequest Request Object
type ShowEngineInstanceExtendProductInfoRequest struct {

	// 消息引擎。
	Engine ShowEngineInstanceExtendProductInfoRequestEngine `json:"engine"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 产品的类型。 - advanced: 专享版
	Type ShowEngineInstanceExtendProductInfoRequestType `json:"type"`
}

func (o ShowEngineInstanceExtendProductInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowEngineInstanceExtendProductInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowEngineInstanceExtendProductInfoRequest", string(data)}, " ")
}

type ShowEngineInstanceExtendProductInfoRequestEngine struct {
	value string
}

type ShowEngineInstanceExtendProductInfoRequestEngineEnum struct {
	KAFKA ShowEngineInstanceExtendProductInfoRequestEngine
}

func GetShowEngineInstanceExtendProductInfoRequestEngineEnum() ShowEngineInstanceExtendProductInfoRequestEngineEnum {
	return ShowEngineInstanceExtendProductInfoRequestEngineEnum{
		KAFKA: ShowEngineInstanceExtendProductInfoRequestEngine{
			value: "kafka",
		},
	}
}

func (c ShowEngineInstanceExtendProductInfoRequestEngine) Value() string {
	return c.value
}

func (c ShowEngineInstanceExtendProductInfoRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowEngineInstanceExtendProductInfoRequestEngine) UnmarshalJSON(b []byte) error {
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

type ShowEngineInstanceExtendProductInfoRequestType struct {
	value string
}

type ShowEngineInstanceExtendProductInfoRequestTypeEnum struct {
	ADVANCED ShowEngineInstanceExtendProductInfoRequestType
}

func GetShowEngineInstanceExtendProductInfoRequestTypeEnum() ShowEngineInstanceExtendProductInfoRequestTypeEnum {
	return ShowEngineInstanceExtendProductInfoRequestTypeEnum{
		ADVANCED: ShowEngineInstanceExtendProductInfoRequestType{
			value: "advanced",
		},
	}
}

func (c ShowEngineInstanceExtendProductInfoRequestType) Value() string {
	return c.value
}

func (c ShowEngineInstanceExtendProductInfoRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowEngineInstanceExtendProductInfoRequestType) UnmarshalJSON(b []byte) error {
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
