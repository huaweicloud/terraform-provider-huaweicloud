package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowInstanceExtendProductInfoRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 产品的类型。 [- advanced: 专享版 - platinum: 铂金版 - dec: 专属云版 - exp: 体验版](tag:hc,hk,hws,hws_hk,ctc,sbc,hk_sbc,cmcc,hws_eu)[- platinum: 铂金版](tag:otc)
	Type ShowInstanceExtendProductInfoRequestType `json:"type"`

	// 消息引擎的类型。当前支持的类型为kafka。
	Engine ShowInstanceExtendProductInfoRequestEngine `json:"engine"`
}

func (o ShowInstanceExtendProductInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceExtendProductInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowInstanceExtendProductInfoRequest", string(data)}, " ")
}

type ShowInstanceExtendProductInfoRequestType struct {
	value string
}

type ShowInstanceExtendProductInfoRequestTypeEnum struct {
	ADVANCED ShowInstanceExtendProductInfoRequestType
	PLATINUM ShowInstanceExtendProductInfoRequestType
	DEC      ShowInstanceExtendProductInfoRequestType
	EXP      ShowInstanceExtendProductInfoRequestType
}

func GetShowInstanceExtendProductInfoRequestTypeEnum() ShowInstanceExtendProductInfoRequestTypeEnum {
	return ShowInstanceExtendProductInfoRequestTypeEnum{
		ADVANCED: ShowInstanceExtendProductInfoRequestType{
			value: "advanced",
		},
		PLATINUM: ShowInstanceExtendProductInfoRequestType{
			value: "platinum",
		},
		DEC: ShowInstanceExtendProductInfoRequestType{
			value: "dec",
		},
		EXP: ShowInstanceExtendProductInfoRequestType{
			value: "exp",
		},
	}
}

func (c ShowInstanceExtendProductInfoRequestType) Value() string {
	return c.value
}

func (c ShowInstanceExtendProductInfoRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowInstanceExtendProductInfoRequestType) UnmarshalJSON(b []byte) error {
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

type ShowInstanceExtendProductInfoRequestEngine struct {
	value string
}

type ShowInstanceExtendProductInfoRequestEngineEnum struct {
	KAFKA ShowInstanceExtendProductInfoRequestEngine
}

func GetShowInstanceExtendProductInfoRequestEngineEnum() ShowInstanceExtendProductInfoRequestEngineEnum {
	return ShowInstanceExtendProductInfoRequestEngineEnum{
		KAFKA: ShowInstanceExtendProductInfoRequestEngine{
			value: "kafka",
		},
	}
}

func (c ShowInstanceExtendProductInfoRequestEngine) Value() string {
	return c.value
}

func (c ShowInstanceExtendProductInfoRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowInstanceExtendProductInfoRequestEngine) UnmarshalJSON(b []byte) error {
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
