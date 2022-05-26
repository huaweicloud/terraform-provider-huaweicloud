package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// oidc详细信息
type UpdateOpenIdConnectConfig struct {

	// 访问方式: program_console: 支持编程访问和管理控制台访问方式; program: 支持编程访问方式
	AccessMode *string `json:"access_mode,omitempty"`

	// OpenID Connect身份提供商标识，对应ID token 中 iss
	IdpUrl *string `json:"idp_url,omitempty"`

	// 在OpenID Connect身份提供商注册的客户端ID
	ClientId *string `json:"client_id,omitempty"`

	// OpenID Connect身份提供商授权地址；编程访问和管理控制台访问方式值不可为空，编程访问方式值可为空
	AuthorizationEndpoint *string `json:"authorization_endpoint,omitempty"`

	// 授权请求信息范围，编程访问和管理控制台访问方式必选，编程访问方式不可选，可选值：openid 、email、profile，IDP自定义scope，字符集a-zA-Z_0-9 ，1-10个可选值组合空格分割，至少包括openid，顺序无关，总长度最长255字符，例如：\"openid\"、\"openid email\"、\"openid profile\" 、\"openid email profile\"
	Scope *string `json:"scope,omitempty"`

	// 授权请求返回的类型；值为id_token ；编程访问和管理控制台访问方式值不可为空，编程访问方式值可为空
	ResponseType *UpdateOpenIdConnectConfigResponseType `json:"response_type,omitempty"`

	// 授权请求返回方式，可选值 form_post 或 fragment ；编程访问和管理控制台访问方式值为可选值，编程访问方式值可为空
	ResponseMode *UpdateOpenIdConnectConfigResponseMode `json:"response_mode,omitempty"`

	// OpenID Connect身份提供商ID Token签名的公钥
	SigningKey *string `json:"signing_key,omitempty"`
}

func (o UpdateOpenIdConnectConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateOpenIdConnectConfig struct{}"
	}

	return strings.Join([]string{"UpdateOpenIdConnectConfig", string(data)}, " ")
}

type UpdateOpenIdConnectConfigResponseType struct {
	value string
}

type UpdateOpenIdConnectConfigResponseTypeEnum struct {
	ID_TOKEN UpdateOpenIdConnectConfigResponseType
}

func GetUpdateOpenIdConnectConfigResponseTypeEnum() UpdateOpenIdConnectConfigResponseTypeEnum {
	return UpdateOpenIdConnectConfigResponseTypeEnum{
		ID_TOKEN: UpdateOpenIdConnectConfigResponseType{
			value: "id_token",
		},
	}
}

func (c UpdateOpenIdConnectConfigResponseType) Value() string {
	return c.value
}

func (c UpdateOpenIdConnectConfigResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateOpenIdConnectConfigResponseType) UnmarshalJSON(b []byte) error {
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

type UpdateOpenIdConnectConfigResponseMode struct {
	value string
}

type UpdateOpenIdConnectConfigResponseModeEnum struct {
	FRAGMENT  UpdateOpenIdConnectConfigResponseMode
	FORM_POST UpdateOpenIdConnectConfigResponseMode
}

func GetUpdateOpenIdConnectConfigResponseModeEnum() UpdateOpenIdConnectConfigResponseModeEnum {
	return UpdateOpenIdConnectConfigResponseModeEnum{
		FRAGMENT: UpdateOpenIdConnectConfigResponseMode{
			value: "fragment",
		},
		FORM_POST: UpdateOpenIdConnectConfigResponseMode{
			value: "form_post",
		},
	}
}

func (c UpdateOpenIdConnectConfigResponseMode) Value() string {
	return c.value
}

func (c UpdateOpenIdConnectConfigResponseMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateOpenIdConnectConfigResponseMode) UnmarshalJSON(b []byte) error {
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
