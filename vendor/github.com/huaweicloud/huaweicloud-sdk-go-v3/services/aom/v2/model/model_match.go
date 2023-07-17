package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Match 静默规则的匹配条件
type Match struct {

	// 指定按照Metadata中的key进行匹配
	Key string `json:"key"`

	// 指定匹配的方式：EXIST:存在，REGEX:正则，EQUALS:等于
	Operate MatchOperate `json:"operate"`

	// 要匹配的key对应的value，当operate为存在时，此值为空
	Value *[]string `json:"value,omitempty"`
}

func (o Match) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Match struct{}"
	}

	return strings.Join([]string{"Match", string(data)}, " ")
}

type MatchOperate struct {
	value string
}

type MatchOperateEnum struct {
	EQUALS MatchOperate
	REGEX  MatchOperate
	EXIST  MatchOperate
}

func GetMatchOperateEnum() MatchOperateEnum {
	return MatchOperateEnum{
		EQUALS: MatchOperate{
			value: "EQUALS",
		},
		REGEX: MatchOperate{
			value: "REGEX",
		},
		EXIST: MatchOperate{
			value: "EXIST",
		},
	}
}

func (c MatchOperate) Value() string {
	return c.value
}

func (c MatchOperate) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MatchOperate) UnmarshalJSON(b []byte) error {
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
