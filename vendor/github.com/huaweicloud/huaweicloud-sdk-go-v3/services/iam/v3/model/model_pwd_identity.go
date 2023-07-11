package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// PwdIdentity
type PwdIdentity struct {

	// 认证方法，该字段内容为[\"password\"]。
	Methods []PwdIdentityMethods `json:"methods"`

	Password *PwdPassword `json:"password"`
}

func (o PwdIdentity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PwdIdentity struct{}"
	}

	return strings.Join([]string{"PwdIdentity", string(data)}, " ")
}

type PwdIdentityMethods struct {
	value string
}

type PwdIdentityMethodsEnum struct {
	PASSWORD PwdIdentityMethods
}

func GetPwdIdentityMethodsEnum() PwdIdentityMethodsEnum {
	return PwdIdentityMethodsEnum{
		PASSWORD: PwdIdentityMethods{
			value: "password",
		},
	}
}

func (c PwdIdentityMethods) Value() string {
	return c.value
}

func (c PwdIdentityMethods) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PwdIdentityMethods) UnmarshalJSON(b []byte) error {
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
