package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AddSubtitle struct {

	// 字幕类型，字幕封装当前仅支持VTT
	Type AddSubtitleType `json:"type"`

	// 字幕语言
	Language string `json:"language"`

	ObsInfo *ObsInfo `json:"obs_info"`
}

func (o AddSubtitle) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSubtitle struct{}"
	}

	return strings.Join([]string{"AddSubtitle", string(data)}, " ")
}

type AddSubtitleType struct {
	value string
}

type AddSubtitleTypeEnum struct {
	VTT AddSubtitleType
}

func GetAddSubtitleTypeEnum() AddSubtitleTypeEnum {
	return AddSubtitleTypeEnum{
		VTT: AddSubtitleType{
			value: "VTT",
		},
	}
}

func (c AddSubtitleType) Value() string {
	return c.value
}

func (c AddSubtitleType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AddSubtitleType) UnmarshalJSON(b []byte) error {
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
