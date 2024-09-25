package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// PolicyStatement
type PolicyStatement struct {

	// 授权项，指对资源的具体操作权限，不超过100个。 > - 格式为：服务名:资源类型:操作，例：vpc:ports:create。 > - 服务名为产品名称，例如ecs、evs和vpc等，服务名仅支持小写。 资源类型和操作没有大小写，要求支持通配符号*，无需罗列全部授权项。 > - 当自定义策略为委托自定义策略时，该字段值为：``` \"Action\": [\"iam:agencies:assume\"]```。
	Action []string `json:"Action"`

	// 作用。包含两种：允许（Allow）和拒绝（Deny），既有Allow又有Deny的授权语句时，遵循Deny优先的原则。
	Effect PolicyStatementEffect `json:"Effect"`

	// 限制条件。不超过10个。
	Condition *interface{} `json:"Condition,omitempty"`

	// 资源。数组长度不超过10，每个字符串长度不超过128，规则如下： > - 可填 * 的五段式：<service-name>:<region>:<account-id>:<resource-type>:<resource-path>，例：\"obs:*:*:bucket:*\"。 > - region字段为*或用户可访问的region。service必须存在且resource属于对应service。 > - 当该自定义策略为委托自定义策略时，该字段类型为Object，值为：```\"Resource\": {\"uri\": [\"/iam/agencies/07805acaba800fdd4fbdc00b8f888c7c\"]}```。
	Resource *interface{} `json:"Resource,omitempty"`
}

func (o PolicyStatement) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyStatement struct{}"
	}

	return strings.Join([]string{"PolicyStatement", string(data)}, " ")
}

type PolicyStatementEffect struct {
	value string
}

type PolicyStatementEffectEnum struct {
	ALLOW PolicyStatementEffect
	DENY  PolicyStatementEffect
}

func GetPolicyStatementEffectEnum() PolicyStatementEffectEnum {
	return PolicyStatementEffectEnum{
		ALLOW: PolicyStatementEffect{
			value: "Allow",
		},
		DENY: PolicyStatementEffect{
			value: "Deny",
		},
	}
}

func (c PolicyStatementEffect) Value() string {
	return c.value
}

func (c PolicyStatementEffect) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PolicyStatementEffect) UnmarshalJSON(b []byte) error {
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
