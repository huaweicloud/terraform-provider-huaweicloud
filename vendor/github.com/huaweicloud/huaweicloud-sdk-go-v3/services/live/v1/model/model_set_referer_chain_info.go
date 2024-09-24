package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type SetRefererChainInfo struct {

	// 直播域名
	Domain string `json:"domain"`

	// referer防盗链开关：true表示开启；false表示关闭
	GuardSwitch SetRefererChainInfoGuardSwitch `json:"guard_switch"`

	// 是否包含referer头域：true表示包含；false表示不包含；guard_switch为true则必填
	RefererConfigEmpty *SetRefererChainInfoRefererConfigEmpty `json:"referer_config_empty,omitempty"`

	// 是否为referer白名单：true表示白名单；false表示黑名单；guard_switch为true则必填
	RefererWhiteList *SetRefererChainInfoRefererWhiteList `json:"referer_white_list,omitempty"`

	// 域名列表，域名为正则表达式；最多支持配置100个域名，以英文“;”进行分隔；guard_switch为true则必填
	RefererAuthList *[]string `json:"referer_auth_list,omitempty"`
}

func (o SetRefererChainInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetRefererChainInfo struct{}"
	}

	return strings.Join([]string{"SetRefererChainInfo", string(data)}, " ")
}

type SetRefererChainInfoGuardSwitch struct {
	value string
}

type SetRefererChainInfoGuardSwitchEnum struct {
	TRUE  SetRefererChainInfoGuardSwitch
	FALSE SetRefererChainInfoGuardSwitch
}

func GetSetRefererChainInfoGuardSwitchEnum() SetRefererChainInfoGuardSwitchEnum {
	return SetRefererChainInfoGuardSwitchEnum{
		TRUE: SetRefererChainInfoGuardSwitch{
			value: "true",
		},
		FALSE: SetRefererChainInfoGuardSwitch{
			value: "false",
		},
	}
}

func (c SetRefererChainInfoGuardSwitch) Value() string {
	return c.value
}

func (c SetRefererChainInfoGuardSwitch) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SetRefererChainInfoGuardSwitch) UnmarshalJSON(b []byte) error {
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

type SetRefererChainInfoRefererConfigEmpty struct {
	value string
}

type SetRefererChainInfoRefererConfigEmptyEnum struct {
	TRUE  SetRefererChainInfoRefererConfigEmpty
	FALSE SetRefererChainInfoRefererConfigEmpty
}

func GetSetRefererChainInfoRefererConfigEmptyEnum() SetRefererChainInfoRefererConfigEmptyEnum {
	return SetRefererChainInfoRefererConfigEmptyEnum{
		TRUE: SetRefererChainInfoRefererConfigEmpty{
			value: "true",
		},
		FALSE: SetRefererChainInfoRefererConfigEmpty{
			value: "false",
		},
	}
}

func (c SetRefererChainInfoRefererConfigEmpty) Value() string {
	return c.value
}

func (c SetRefererChainInfoRefererConfigEmpty) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SetRefererChainInfoRefererConfigEmpty) UnmarshalJSON(b []byte) error {
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

type SetRefererChainInfoRefererWhiteList struct {
	value string
}

type SetRefererChainInfoRefererWhiteListEnum struct {
	TRUE  SetRefererChainInfoRefererWhiteList
	FALSE SetRefererChainInfoRefererWhiteList
}

func GetSetRefererChainInfoRefererWhiteListEnum() SetRefererChainInfoRefererWhiteListEnum {
	return SetRefererChainInfoRefererWhiteListEnum{
		TRUE: SetRefererChainInfoRefererWhiteList{
			value: "true",
		},
		FALSE: SetRefererChainInfoRefererWhiteList{
			value: "false",
		},
	}
}

func (c SetRefererChainInfoRefererWhiteList) Value() string {
	return c.value
}

func (c SetRefererChainInfoRefererWhiteList) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SetRefererChainInfoRefererWhiteList) UnmarshalJSON(b []byte) error {
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
