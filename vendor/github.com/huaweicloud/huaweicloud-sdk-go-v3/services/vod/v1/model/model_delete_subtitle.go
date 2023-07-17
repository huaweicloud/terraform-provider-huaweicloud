package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type DeleteSubtitle struct {

	// 字幕类型，字幕封装当前仅支持VTT
	Type DeleteSubtitleType `json:"type"`

	// 字幕语言
	Language string `json:"language"`
}

func (o DeleteSubtitle) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSubtitle struct{}"
	}

	return strings.Join([]string{"DeleteSubtitle", string(data)}, " ")
}

type DeleteSubtitleType struct {
	value string
}

type DeleteSubtitleTypeEnum struct {
	VTT DeleteSubtitleType
}

func GetDeleteSubtitleTypeEnum() DeleteSubtitleTypeEnum {
	return DeleteSubtitleTypeEnum{
		VTT: DeleteSubtitleType{
			value: "VTT",
		},
	}
}

func (c DeleteSubtitleType) Value() string {
	return c.value
}

func (c DeleteSubtitleType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteSubtitleType) UnmarshalJSON(b []byte) error {
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
