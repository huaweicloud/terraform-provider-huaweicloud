package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

//
type ModifyTransTemplateGroup struct {

	// 模板组名称。
	GroupId string `json:"group_id"`

	// 模板组名称。
	Name string `json:"name"`

	// 是否设置默认。
	Status *ModifyTransTemplateGroupStatus `json:"status,omitempty"`

	// 是否自动加密。  取值如下： - 0：表示不加密。 - 1：表示需要加密。  默认值：0。  加密与转码必须要一起进行，当需要加密时，转码参数不能为空，且转码输出格式必须要为HLS。
	AutoEncrypt *int32 `json:"auto_encrypt,omitempty"`

	// 画质配置信息列表。
	QualityInfoList *[]QualityInfo `json:"quality_info_list,omitempty"`

	// 绑定的水印模板组ID数组。
	WatermarkTemplateIds *[]string `json:"watermark_template_ids,omitempty"`

	// 模板介绍。
	Description *string `json:"description,omitempty"`

	Common *Common `json:"common,omitempty"`
}

func (o ModifyTransTemplateGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyTransTemplateGroup struct{}"
	}

	return strings.Join([]string{"ModifyTransTemplateGroup", string(data)}, " ")
}

type ModifyTransTemplateGroupStatus struct {
	value string
}

type ModifyTransTemplateGroupStatusEnum struct {
	E_1 ModifyTransTemplateGroupStatus
	E_0 ModifyTransTemplateGroupStatus
}

func GetModifyTransTemplateGroupStatusEnum() ModifyTransTemplateGroupStatusEnum {
	return ModifyTransTemplateGroupStatusEnum{
		E_1: ModifyTransTemplateGroupStatus{
			value: "1",
		},
		E_0: ModifyTransTemplateGroupStatus{
			value: "0",
		},
	}
}

func (c ModifyTransTemplateGroupStatus) Value() string {
	return c.value
}

func (c ModifyTransTemplateGroupStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ModifyTransTemplateGroupStatus) UnmarshalJSON(b []byte) error {
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
