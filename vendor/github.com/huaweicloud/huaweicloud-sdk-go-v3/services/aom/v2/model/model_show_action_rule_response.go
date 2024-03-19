package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowActionRuleResponse Response Object
type ShowActionRuleResponse struct {

	// 规则名称 只含有汉字数字、字母、下划线，不能以下划线等特殊符号开头和结尾，长度为 1 - 100
	RuleName *string `json:"rule_name,omitempty"`

	// 项目ID
	ProjectId *string `json:"project_id,omitempty"`

	// 子账号名称
	UserName *string `json:"user_name,omitempty"`

	// 规则描述。规则描述长度为0到1024个字符，并且只能是数字、字母、特殊字符（_*）、空格和中文组成，不能以下划线开头和结尾。
	Desc *string `json:"desc,omitempty"`

	// 规则类型。\"1\"：通知，\"2\"：用户
	Type *ShowActionRuleResponseType `json:"type,omitempty"`

	// 消息模板
	NotificationTemplate *string `json:"notification_template,omitempty"`

	// 创建时间
	CreateTime *int64 `json:"create_time,omitempty"`

	// 修改时间
	UpdateTime *int64 `json:"update_time,omitempty"`

	// 时区
	TimeZone *string `json:"time_zone,omitempty"`

	// SMN主题信息，不能大于5
	SmnTopics      *[]SmnTopics `json:"smn_topics,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ShowActionRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowActionRuleResponse struct{}"
	}

	return strings.Join([]string{"ShowActionRuleResponse", string(data)}, " ")
}

type ShowActionRuleResponseType struct {
	value string
}

type ShowActionRuleResponseTypeEnum struct {
	E_1 ShowActionRuleResponseType
	E_2 ShowActionRuleResponseType
}

func GetShowActionRuleResponseTypeEnum() ShowActionRuleResponseTypeEnum {
	return ShowActionRuleResponseTypeEnum{
		E_1: ShowActionRuleResponseType{
			value: "1",
		},
		E_2: ShowActionRuleResponseType{
			value: "2",
		},
	}
}

func (c ShowActionRuleResponseType) Value() string {
	return c.value
}

func (c ShowActionRuleResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowActionRuleResponseType) UnmarshalJSON(b []byte) error {
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
