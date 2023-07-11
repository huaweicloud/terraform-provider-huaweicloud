package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type PolicyEntity struct {

	// 是否为创建topic时所选择的用户。
	Owner *bool `json:"owner,omitempty"`

	// 用户名。
	UserName *string `json:"user_name,omitempty"`

	// 权限类型。 - all：拥有发布、订阅权限; - pub：拥有发布权限; - sub：拥有订阅权限。
	AccessPolicy *PolicyEntityAccessPolicy `json:"access_policy,omitempty"`
}

func (o PolicyEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyEntity struct{}"
	}

	return strings.Join([]string{"PolicyEntity", string(data)}, " ")
}

type PolicyEntityAccessPolicy struct {
	value string
}

type PolicyEntityAccessPolicyEnum struct {
	ALL PolicyEntityAccessPolicy
	PUB PolicyEntityAccessPolicy
	SUB PolicyEntityAccessPolicy
}

func GetPolicyEntityAccessPolicyEnum() PolicyEntityAccessPolicyEnum {
	return PolicyEntityAccessPolicyEnum{
		ALL: PolicyEntityAccessPolicy{
			value: "all",
		},
		PUB: PolicyEntityAccessPolicy{
			value: "pub",
		},
		SUB: PolicyEntityAccessPolicy{
			value: "sub",
		},
	}
}

func (c PolicyEntityAccessPolicy) Value() string {
	return c.value
}

func (c PolicyEntityAccessPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PolicyEntityAccessPolicy) UnmarshalJSON(b []byte) error {
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
