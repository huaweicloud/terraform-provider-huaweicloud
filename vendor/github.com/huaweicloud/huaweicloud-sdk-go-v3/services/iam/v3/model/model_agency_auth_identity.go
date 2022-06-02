package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

//
type AgencyAuthIdentity struct {

	// 认证方法，该字段内容为[\"assume_role\"]。
	Methods []AgencyAuthIdentityMethods `json:"methods"`

	AssumeRole *IdentityAssumerole `json:"assume_role"`

	Policy *ServicePolicy `json:"policy,omitempty"`
}

func (o AgencyAuthIdentity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyAuthIdentity struct{}"
	}

	return strings.Join([]string{"AgencyAuthIdentity", string(data)}, " ")
}

type AgencyAuthIdentityMethods struct {
	value string
}

type AgencyAuthIdentityMethodsEnum struct {
	ASSUME_ROLE AgencyAuthIdentityMethods
}

func GetAgencyAuthIdentityMethodsEnum() AgencyAuthIdentityMethodsEnum {
	return AgencyAuthIdentityMethodsEnum{
		ASSUME_ROLE: AgencyAuthIdentityMethods{
			value: "assume_role",
		},
	}
}

func (c AgencyAuthIdentityMethods) Value() string {
	return c.value
}

func (c AgencyAuthIdentityMethods) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AgencyAuthIdentityMethods) UnmarshalJSON(b []byte) error {
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
