package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// oidc配置详细信息
type OpenIdConnectConfig struct {

	// 访问方式: program_console: 支持编程访问和管理控制台访问方式; program: 支持编程访问方式
	AccessMode string `json:"access_mode"`

	// OpenID Connect身份提供商标识, 对应ID token 中 iss
	IdpUrl string `json:"idp_url"`

	// 在OpenID Connect身份提供商注册的客户端ID
	ClientId string `json:"client_id"`

	// OpenID Connect身份提供商授权地址; 编程访问和管理控制台访问方式值不为空，编程访问方式值可为空
	AuthorizationEndpoint string `json:"authorization_endpoint"`

	// 授权请求信息范围，编程访问和管理控制台访问方式必选，编程访问方式不可选，可选值：openid 、email、profile，IDP自定义scope，字符集a-zA-Z_0-9 ，1-10个可选值组合空格分割，至少包括openid，顺序无关，总长度最长255字符，例如：\"openid\"、\"openid email\"、\"openid profile\" 、\"openid email profile\"
	Scope string `json:"scope"`

	// 授权请求返回的类型；id_token ；编程访问和管理控制台访问方式值不为空，编程访问方式值可为空
	ResponseType OpenIdConnectConfigResponseType `json:"response_type"`

	// 授权请求返回方式， form_post 或 fragment ；编程访问和管理控制台访问方式值不为空，编程访问方式值可为空
	ResponseMode OpenIdConnectConfigResponseMode `json:"response_mode"`

	// OpenID Connect身份提供商ID Token签名的公钥
	SigningKey string `json:"signing_key"`
}

func (o OpenIdConnectConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OpenIdConnectConfig struct{}"
	}

	return strings.Join([]string{"OpenIdConnectConfig", string(data)}, " ")
}

type OpenIdConnectConfigResponseType struct {
	value string
}

type OpenIdConnectConfigResponseTypeEnum struct {
	ID_TOKEN OpenIdConnectConfigResponseType
}

func GetOpenIdConnectConfigResponseTypeEnum() OpenIdConnectConfigResponseTypeEnum {
	return OpenIdConnectConfigResponseTypeEnum{
		ID_TOKEN: OpenIdConnectConfigResponseType{
			value: "id_token",
		},
	}
}

func (c OpenIdConnectConfigResponseType) Value() string {
	return c.value
}

func (c OpenIdConnectConfigResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OpenIdConnectConfigResponseType) UnmarshalJSON(b []byte) error {
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

type OpenIdConnectConfigResponseMode struct {
	value string
}

type OpenIdConnectConfigResponseModeEnum struct {
	FRAGMENT  OpenIdConnectConfigResponseMode
	FORM_POST OpenIdConnectConfigResponseMode
}

func GetOpenIdConnectConfigResponseModeEnum() OpenIdConnectConfigResponseModeEnum {
	return OpenIdConnectConfigResponseModeEnum{
		FRAGMENT: OpenIdConnectConfigResponseMode{
			value: "fragment",
		},
		FORM_POST: OpenIdConnectConfigResponseMode{
			value: "form_post",
		},
	}
}

func (c OpenIdConnectConfigResponseMode) Value() string {
	return c.value
}

func (c OpenIdConnectConfigResponseMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OpenIdConnectConfigResponseMode) UnmarshalJSON(b []byte) error {
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
