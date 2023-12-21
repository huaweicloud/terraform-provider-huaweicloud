package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CallbackUrl struct {

	// 回调地址
	Url string `json:"url"`

	// 鉴权密钥
	AuthSignKey *string `json:"auth_sign_key,omitempty"`

	// 接收回调通知服务器所在区域。 包含如下取值： - mainland_china：中国大陆区域。 - outside_mainland_china：中国大陆以外区域。
	CallBackArea *CallbackUrlCallBackArea `json:"call_back_area,omitempty"`
}

func (o CallbackUrl) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CallbackUrl struct{}"
	}

	return strings.Join([]string{"CallbackUrl", string(data)}, " ")
}

type CallbackUrlCallBackArea struct {
	value string
}

type CallbackUrlCallBackAreaEnum struct {
	MAINLAND_CHINA         CallbackUrlCallBackArea
	OUTSIDE_MAINLAND_CHINA CallbackUrlCallBackArea
}

func GetCallbackUrlCallBackAreaEnum() CallbackUrlCallBackAreaEnum {
	return CallbackUrlCallBackAreaEnum{
		MAINLAND_CHINA: CallbackUrlCallBackArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: CallbackUrlCallBackArea{
			value: "outside_mainland_china",
		},
	}
}

func (c CallbackUrlCallBackArea) Value() string {
	return c.value
}

func (c CallbackUrlCallBackArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CallbackUrlCallBackArea) UnmarshalJSON(b []byte) error {
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
