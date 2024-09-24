package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowPullSourcesConfigResponse Response Object
type ShowPullSourcesConfigResponse struct {

	// 播放域名
	PlayDomain *string `json:"play_domain,omitempty"`

	// 回源方式。 - domain: 回源客户源站，源站地址是域名格式。回源域名，可配置多个，如果回源失败，将按照配置顺序进行轮循。 - ipaddr: 回源客户源站，源站地址是IP格式。回源IP，可配置多个，如果回源失败，将按照配置顺序进行轮循。同时，最多可以配置一个回源域名，如果配置，回源时httpflv HOST头填该域名，RTMP tcurl字段填该域名，否则按当前IP作为HOST。 - huawei: 回源华为源站，域名创建后的默认值。
	SourceType *ShowPullSourcesConfigResponseSourceType `json:"source_type,omitempty"`

	// 回源域名列表，最多可配置10个。 - 当回源方式是“domain”时，此参数必选，域名配置多个时，如果回源失败，将按照配置顺序进行轮循。 - 当回源方式是“ipaddr”时，最多可以配置一个回源域名，如果配置，回源时httpflv HOST头填该域名，RTMP tcurl 字段填该域名，否则按当前IP作为HOST。
	Sources *[]string `json:"sources,omitempty"`

	// 回源IP地址列表，最多可配置10个。当回源方式是“ipaddr”时，此参数必选，IP配置多个时，如果回源失败，将按照配置顺序进行轮循。
	SourcesIp *[]string `json:"sources_ip,omitempty"`

	// 回源端口。
	SourcePort *int32 `json:"source_port,omitempty"`

	// 回源协议，回源方式非“huawei”时必选。
	Scheme *ShowPullSourcesConfigResponseScheme `json:"scheme,omitempty"`

	// 回源客户源站时在URL携带的参数。
	AdditionalArgs map[string]string `json:"additional_args,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ShowPullSourcesConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPullSourcesConfigResponse struct{}"
	}

	return strings.Join([]string{"ShowPullSourcesConfigResponse", string(data)}, " ")
}

type ShowPullSourcesConfigResponseSourceType struct {
	value string
}

type ShowPullSourcesConfigResponseSourceTypeEnum struct {
	DOMAIN ShowPullSourcesConfigResponseSourceType
	IPADDR ShowPullSourcesConfigResponseSourceType
	HUAWEI ShowPullSourcesConfigResponseSourceType
}

func GetShowPullSourcesConfigResponseSourceTypeEnum() ShowPullSourcesConfigResponseSourceTypeEnum {
	return ShowPullSourcesConfigResponseSourceTypeEnum{
		DOMAIN: ShowPullSourcesConfigResponseSourceType{
			value: "domain",
		},
		IPADDR: ShowPullSourcesConfigResponseSourceType{
			value: "ipaddr",
		},
		HUAWEI: ShowPullSourcesConfigResponseSourceType{
			value: "huawei",
		},
	}
}

func (c ShowPullSourcesConfigResponseSourceType) Value() string {
	return c.value
}

func (c ShowPullSourcesConfigResponseSourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowPullSourcesConfigResponseSourceType) UnmarshalJSON(b []byte) error {
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

type ShowPullSourcesConfigResponseScheme struct {
	value string
}

type ShowPullSourcesConfigResponseSchemeEnum struct {
	HTTP ShowPullSourcesConfigResponseScheme
	RTMP ShowPullSourcesConfigResponseScheme
}

func GetShowPullSourcesConfigResponseSchemeEnum() ShowPullSourcesConfigResponseSchemeEnum {
	return ShowPullSourcesConfigResponseSchemeEnum{
		HTTP: ShowPullSourcesConfigResponseScheme{
			value: "http",
		},
		RTMP: ShowPullSourcesConfigResponseScheme{
			value: "rtmp",
		},
	}
}

func (c ShowPullSourcesConfigResponseScheme) Value() string {
	return c.value
}

func (c ShowPullSourcesConfigResponseScheme) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowPullSourcesConfigResponseScheme) UnmarshalJSON(b []byte) error {
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
