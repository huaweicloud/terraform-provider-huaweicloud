package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AgencyPolicyStatement
type AgencyPolicyStatement struct {

	// 授权项，指对资源的具体操作权限。 > - 当自定义策略为委托自定义策略时，该字段值为：``` \"Action\": [\"iam:agencies:assume\"]```。
	Action []AgencyPolicyStatementAction `json:"Action"`

	// 作用。包含两种：允许（Allow）和拒绝（Deny），既有Allow又有Deny的授权语句时，遵循Deny优先的原则。
	Effect AgencyPolicyStatementEffect `json:"Effect"`

	Resource *AgencyPolicyResource `json:"Resource,omitempty"`
}

func (o AgencyPolicyStatement) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyPolicyStatement struct{}"
	}

	return strings.Join([]string{"AgencyPolicyStatement", string(data)}, " ")
}

type AgencyPolicyStatementAction struct {
	value string
}

type AgencyPolicyStatementActionEnum struct {
	IAMAGENCIESASSUME AgencyPolicyStatementAction
}

func GetAgencyPolicyStatementActionEnum() AgencyPolicyStatementActionEnum {
	return AgencyPolicyStatementActionEnum{
		IAMAGENCIESASSUME: AgencyPolicyStatementAction{
			value: "iam:agencies:assume",
		},
	}
}

func (c AgencyPolicyStatementAction) Value() string {
	return c.value
}

func (c AgencyPolicyStatementAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AgencyPolicyStatementAction) UnmarshalJSON(b []byte) error {
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

type AgencyPolicyStatementEffect struct {
	value string
}

type AgencyPolicyStatementEffectEnum struct {
	ALLOW AgencyPolicyStatementEffect
	DENY  AgencyPolicyStatementEffect
}

func GetAgencyPolicyStatementEffectEnum() AgencyPolicyStatementEffectEnum {
	return AgencyPolicyStatementEffectEnum{
		ALLOW: AgencyPolicyStatementEffect{
			value: "Allow",
		},
		DENY: AgencyPolicyStatementEffect{
			value: "Deny",
		},
	}
}

func (c AgencyPolicyStatementEffect) Value() string {
	return c.value
}

func (c AgencyPolicyStatementEffect) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AgencyPolicyStatementEffect) UnmarshalJSON(b []byte) error {
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
