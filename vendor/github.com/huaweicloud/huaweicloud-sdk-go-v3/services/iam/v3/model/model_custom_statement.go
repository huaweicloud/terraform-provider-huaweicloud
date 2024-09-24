package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CustomStatement
type CustomStatement struct {

	// 授权项，指对资源的具体操作权限。 > - 格式为：服务名:资源类型:操作，例：vpc:ports:create。 > - 服务名为产品名称，例如ecs、evs和vpc等，服务名仅支持小写。 资源类型和操作没有大小写，要求支持通配符号*，无需罗列全部授权项。
	Action []string `json:"Action"`

	// 作用。包含两种：允许（Allow）和拒绝（Deny），既有Allow又有Deny的授权语句时，遵循Deny优先的原则。
	Effect CustomStatementEffect `json:"Effect"`

	Condition map[string]map[string][]string `json:"Condition,omitempty"`

	// 资源。规则如下： > - 可填 * 的五段式：<service-name>:<region>:<account-id>:<resource-type>:<resource-path>，例：\"obs:*:*:bucket:*\"。 > - region字段为*或用户可访问的region。service必须存在且resource属于对应service。
	Resource *interface{} `json:"Resource,omitempty"`
}

func (o CustomStatement) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CustomStatement struct{}"
	}

	return strings.Join([]string{"CustomStatement", string(data)}, " ")
}

type CustomStatementEffect struct {
	value string
}

type CustomStatementEffectEnum struct {
	ALLOW CustomStatementEffect
	DENY  CustomStatementEffect
}

func GetCustomStatementEffectEnum() CustomStatementEffectEnum {
	return CustomStatementEffectEnum{
		ALLOW: CustomStatementEffect{
			value: "Allow",
		},
		DENY: CustomStatementEffect{
			value: "Deny",
		},
	}
}

func (c CustomStatementEffect) Value() string {
	return c.value
}

func (c CustomStatementEffect) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CustomStatementEffect) UnmarshalJSON(b []byte) error {
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
