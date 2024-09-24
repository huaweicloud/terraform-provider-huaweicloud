package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowRefererChainResponse Response Object
type ShowRefererChainResponse struct {

	// 直播域名
	Domain *string `json:"domain,omitempty"`

	// referer防盗链开关：true表示开启；false表示关闭
	GuardSwitch *ShowRefererChainResponseGuardSwitch `json:"guard_switch,omitempty"`

	// 是否包含referer头域：true表示包含；false表示不包含；guard_switch为true则必填
	RefererConfigEmpty *ShowRefererChainResponseRefererConfigEmpty `json:"referer_config_empty,omitempty"`

	// 是否为referer白名单：true表示白名单；false表示黑名单；guard_switch为true则必填
	RefererWhiteList *ShowRefererChainResponseRefererWhiteList `json:"referer_white_list,omitempty"`

	// 域名列表，域名为正则表达式；最多支持配置100个域名，以英文“;”进行分隔；guard_switch为true则必填
	RefererAuthList *[]string `json:"referer_auth_list,omitempty"`
	HttpStatusCode  int       `json:"-"`
}

func (o ShowRefererChainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRefererChainResponse struct{}"
	}

	return strings.Join([]string{"ShowRefererChainResponse", string(data)}, " ")
}

type ShowRefererChainResponseGuardSwitch struct {
	value string
}

type ShowRefererChainResponseGuardSwitchEnum struct {
	TRUE  ShowRefererChainResponseGuardSwitch
	FALSE ShowRefererChainResponseGuardSwitch
}

func GetShowRefererChainResponseGuardSwitchEnum() ShowRefererChainResponseGuardSwitchEnum {
	return ShowRefererChainResponseGuardSwitchEnum{
		TRUE: ShowRefererChainResponseGuardSwitch{
			value: "true",
		},
		FALSE: ShowRefererChainResponseGuardSwitch{
			value: "false",
		},
	}
}

func (c ShowRefererChainResponseGuardSwitch) Value() string {
	return c.value
}

func (c ShowRefererChainResponseGuardSwitch) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowRefererChainResponseGuardSwitch) UnmarshalJSON(b []byte) error {
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

type ShowRefererChainResponseRefererConfigEmpty struct {
	value string
}

type ShowRefererChainResponseRefererConfigEmptyEnum struct {
	TRUE  ShowRefererChainResponseRefererConfigEmpty
	FALSE ShowRefererChainResponseRefererConfigEmpty
}

func GetShowRefererChainResponseRefererConfigEmptyEnum() ShowRefererChainResponseRefererConfigEmptyEnum {
	return ShowRefererChainResponseRefererConfigEmptyEnum{
		TRUE: ShowRefererChainResponseRefererConfigEmpty{
			value: "true",
		},
		FALSE: ShowRefererChainResponseRefererConfigEmpty{
			value: "false",
		},
	}
}

func (c ShowRefererChainResponseRefererConfigEmpty) Value() string {
	return c.value
}

func (c ShowRefererChainResponseRefererConfigEmpty) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowRefererChainResponseRefererConfigEmpty) UnmarshalJSON(b []byte) error {
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

type ShowRefererChainResponseRefererWhiteList struct {
	value string
}

type ShowRefererChainResponseRefererWhiteListEnum struct {
	TRUE  ShowRefererChainResponseRefererWhiteList
	FALSE ShowRefererChainResponseRefererWhiteList
}

func GetShowRefererChainResponseRefererWhiteListEnum() ShowRefererChainResponseRefererWhiteListEnum {
	return ShowRefererChainResponseRefererWhiteListEnum{
		TRUE: ShowRefererChainResponseRefererWhiteList{
			value: "true",
		},
		FALSE: ShowRefererChainResponseRefererWhiteList{
			value: "false",
		},
	}
}

func (c ShowRefererChainResponseRefererWhiteList) Value() string {
	return c.value
}

func (c ShowRefererChainResponseRefererWhiteList) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowRefererChainResponseRefererWhiteList) UnmarshalJSON(b []byte) error {
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
