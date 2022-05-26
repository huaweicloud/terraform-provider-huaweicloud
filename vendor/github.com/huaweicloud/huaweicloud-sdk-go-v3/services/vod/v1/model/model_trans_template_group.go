package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type TransTemplateGroup struct {

	// 模板组名称。
	Name string `json:"name"`

	// 是否设置默认。
	Status *TransTemplateGroupStatus `json:"status,omitempty"`

	// 模板组类型。
	Type TransTemplateGroupType `json:"type"`

	// 是否自动加密。  取值如下： - 0：表示不加密。 - 1：表示需要加密。  默认值：0。  加密与转码必须要一起进行，当需要加密时，转码参数不能为空，且转码输出格式必须要为HLS。
	AutoEncrypt *int32 `json:"auto_encrypt,omitempty"`

	// 画质配置信息列表。
	QualityInfoList *[]QualityInfo `json:"quality_info_list,omitempty"`

	Common *Common `json:"common,omitempty"`

	// 绑定的水印模板组ID数组。
	WatermarkTemplateIds *[]string `json:"watermark_template_ids,omitempty"`

	// 模板介绍。
	Description *string `json:"description,omitempty"`
}

func (o TransTemplateGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TransTemplateGroup struct{}"
	}

	return strings.Join([]string{"TransTemplateGroup", string(data)}, " ")
}

type TransTemplateGroupStatus struct {
	value string
}

type TransTemplateGroupStatusEnum struct {
	E_1 TransTemplateGroupStatus
	E_0 TransTemplateGroupStatus
}

func GetTransTemplateGroupStatusEnum() TransTemplateGroupStatusEnum {
	return TransTemplateGroupStatusEnum{
		E_1: TransTemplateGroupStatus{
			value: "1",
		},
		E_0: TransTemplateGroupStatus{
			value: "0",
		},
	}
}

func (c TransTemplateGroupStatus) Value() string {
	return c.value
}

func (c TransTemplateGroupStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TransTemplateGroupStatus) UnmarshalJSON(b []byte) error {
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

type TransTemplateGroupType struct {
	value string
}

type TransTemplateGroupTypeEnum struct {
	CUSTOM_TEMPLATE_GROUP TransTemplateGroupType
}

func GetTransTemplateGroupTypeEnum() TransTemplateGroupTypeEnum {
	return TransTemplateGroupTypeEnum{
		CUSTOM_TEMPLATE_GROUP: TransTemplateGroupType{
			value: "custom_template_group",
		},
	}
}

func (c TransTemplateGroupType) Value() string {
	return c.value
}

func (c TransTemplateGroupType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TransTemplateGroupType) UnmarshalJSON(b []byte) error {
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
