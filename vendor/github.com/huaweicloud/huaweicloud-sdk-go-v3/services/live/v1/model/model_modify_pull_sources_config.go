package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type ModifyPullSourcesConfig struct {

	// 直播播放域名
	PlayDomain string `json:"play_domain"`

	// 回源方式。  包含如下取值： - domain: 回源客户源站，源站地址是域名格式。回源域名，可配置多个，如果回源失败，将按照配置顺序进行轮循。 - ipaddr: 回源客户源站，源站地址是IP格式。回源IP，可配置多个，如果回源失败，将按照配置顺序进行轮循。同时，最多可以配置一个回源域名，如果配置，回源时httpflv HOST头填该域名，RTMP tcurl字段填该域名，否则按当前IP作为HOST。 - huawei: 回源华为源站，域名创建后的默认值。
	SourceType ModifyPullSourcesConfigSourceType `json:"source_type"`

	// 回源域名列表，最多可配置10个。 - 当回源方式是“domain”时，此参数必选，域名配置多个时，如果回源失败，将按照配置顺序进行轮循。 - 当回源方式是“ipaddr”时，最多可以配置一个回源域名，如果配置，回源时httpflv HOST头填该域名，RTMP tcurl 字段填该域名，否则按当前IP作为HOST。
	Sources *[]string `json:"sources,omitempty"`

	// 回源IP地址列表，最多可配置10个。当回源方式是“ipaddr”时，此参数必选，IP配置多个时，如果回源失败，将按照配置顺序进行轮循。
	SourcesIp *[]string `json:"sources_ip,omitempty"`

	// 回源端口。
	SourcePort *int32 `json:"source_port,omitempty"`

	// 回源协议，回源方式非“huawei”时必选。  包含如下取值： - http - rtmp
	Scheme *ModifyPullSourcesConfigScheme `json:"scheme,omitempty"`

	// 回源客户源站时在URL携带的参数。
	AdditionalArgs map[string]string `json:"additional_args,omitempty"`
}

func (o ModifyPullSourcesConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyPullSourcesConfig struct{}"
	}

	return strings.Join([]string{"ModifyPullSourcesConfig", string(data)}, " ")
}

type ModifyPullSourcesConfigSourceType struct {
	value string
}

type ModifyPullSourcesConfigSourceTypeEnum struct {
	DOMAIN ModifyPullSourcesConfigSourceType
	IPADDR ModifyPullSourcesConfigSourceType
	HUAWEI ModifyPullSourcesConfigSourceType
}

func GetModifyPullSourcesConfigSourceTypeEnum() ModifyPullSourcesConfigSourceTypeEnum {
	return ModifyPullSourcesConfigSourceTypeEnum{
		DOMAIN: ModifyPullSourcesConfigSourceType{
			value: "domain",
		},
		IPADDR: ModifyPullSourcesConfigSourceType{
			value: "ipaddr",
		},
		HUAWEI: ModifyPullSourcesConfigSourceType{
			value: "huawei",
		},
	}
}

func (c ModifyPullSourcesConfigSourceType) Value() string {
	return c.value
}

func (c ModifyPullSourcesConfigSourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ModifyPullSourcesConfigSourceType) UnmarshalJSON(b []byte) error {
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

type ModifyPullSourcesConfigScheme struct {
	value string
}

type ModifyPullSourcesConfigSchemeEnum struct {
	HTTP ModifyPullSourcesConfigScheme
	RTMP ModifyPullSourcesConfigScheme
}

func GetModifyPullSourcesConfigSchemeEnum() ModifyPullSourcesConfigSchemeEnum {
	return ModifyPullSourcesConfigSchemeEnum{
		HTTP: ModifyPullSourcesConfigScheme{
			value: "http",
		},
		RTMP: ModifyPullSourcesConfigScheme{
			value: "rtmp",
		},
	}
}

func (c ModifyPullSourcesConfigScheme) Value() string {
	return c.value
}

func (c ModifyPullSourcesConfigScheme) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ModifyPullSourcesConfigScheme) UnmarshalJSON(b []byte) error {
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
