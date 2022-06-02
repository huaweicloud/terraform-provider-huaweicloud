package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// oidc配置详细信息
type CreateOpenIdConnectConfig struct {

	// 访问方式: program_console: 支持编程访问和管理控制台访问方式; program: 支持编程访问方式
	AccessMode string `json:"access_mode"`

	// OpenID Connect身份提供商标识, 对应ID token 中 iss
	IdpUrl string `json:"idp_url"`

	// 在OpenID Connect身份提供商注册的客户端ID
	ClientId string `json:"client_id"`

	// OpenID Connect身份提供商授权地址;编程访问和管理控制台访问方式必选，编程访问方式不可选
	AuthorizationEndpoint *string `json:"authorization_endpoint,omitempty"`

	// 授权请求信息范围，编程访问和管理控制台访问方式必选，编程访问方式不可选，可选值：openid 、email、profile，IDP自定义scope，字符集a-zA-Z_0-9 ，1-10个可选值组合空格分割，至少包括openid，顺序无关，总长度最长255字符，例如：\"openid\"、\"openid email\"、\"openid profile\" 、\"openid email profile\"
	Scope *string `json:"scope,omitempty"`

	// 授权请求返回的类型；id_token ；编程访问和管理控制台访问方式必选，编程访问方式不可选
	ResponseType *CreateOpenIdConnectConfigResponseType `json:"response_type,omitempty"`

	// 授权请求返回方式， form_post 或 fragment ；编程访问和管理控制台访问方式必选，编程访问方式不可选
	ResponseMode *CreateOpenIdConnectConfigResponseMode `json:"response_mode,omitempty"`

	// OpenID Connect身份提供商ID Token签名的公钥
	SigningKey string `json:"signing_key"`
}

func (o CreateOpenIdConnectConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOpenIdConnectConfig struct{}"
	}

	return strings.Join([]string{"CreateOpenIdConnectConfig", string(data)}, " ")
}

type CreateOpenIdConnectConfigResponseType struct {
	value string
}

type CreateOpenIdConnectConfigResponseTypeEnum struct {
	ID_TOKEN CreateOpenIdConnectConfigResponseType
}

func GetCreateOpenIdConnectConfigResponseTypeEnum() CreateOpenIdConnectConfigResponseTypeEnum {
	return CreateOpenIdConnectConfigResponseTypeEnum{
		ID_TOKEN: CreateOpenIdConnectConfigResponseType{
			value: "id_token",
		},
	}
}

func (c CreateOpenIdConnectConfigResponseType) Value() string {
	return c.value
}

func (c CreateOpenIdConnectConfigResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateOpenIdConnectConfigResponseType) UnmarshalJSON(b []byte) error {
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

type CreateOpenIdConnectConfigResponseMode struct {
	value string
}

type CreateOpenIdConnectConfigResponseModeEnum struct {
	FRAGMENT  CreateOpenIdConnectConfigResponseMode
	FORM_POST CreateOpenIdConnectConfigResponseMode
}

func GetCreateOpenIdConnectConfigResponseModeEnum() CreateOpenIdConnectConfigResponseModeEnum {
	return CreateOpenIdConnectConfigResponseModeEnum{
		FRAGMENT: CreateOpenIdConnectConfigResponseMode{
			value: "fragment",
		},
		FORM_POST: CreateOpenIdConnectConfigResponseMode{
			value: "form_post",
		},
	}
}

func (c CreateOpenIdConnectConfigResponseMode) Value() string {
	return c.value
}

func (c CreateOpenIdConnectConfigResponseMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateOpenIdConnectConfigResponseMode) UnmarshalJSON(b []byte) error {
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
