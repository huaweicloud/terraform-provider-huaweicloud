package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

//
type TokenAuthIdentity struct {

	// 认证方法，该字段内容为[\"token\"]。
	Methods []TokenAuthIdentityMethods `json:"methods"`

	Token *IdentityToken `json:"token,omitempty"`

	Policy *ServicePolicy `json:"policy,omitempty"`
}

func (o TokenAuthIdentity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenAuthIdentity struct{}"
	}

	return strings.Join([]string{"TokenAuthIdentity", string(data)}, " ")
}

type TokenAuthIdentityMethods struct {
	value string
}

type TokenAuthIdentityMethodsEnum struct {
	TOKEN TokenAuthIdentityMethods
}

func GetTokenAuthIdentityMethodsEnum() TokenAuthIdentityMethodsEnum {
	return TokenAuthIdentityMethodsEnum{
		TOKEN: TokenAuthIdentityMethods{
			value: "token",
		},
	}
}

func (c TokenAuthIdentityMethods) Value() string {
	return c.value
}

func (c TokenAuthIdentityMethods) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TokenAuthIdentityMethods) UnmarshalJSON(b []byte) error {
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
