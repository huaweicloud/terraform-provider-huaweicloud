package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Subtitle struct {

	// 字幕id。  取值范围：[1,8]。
	Id int32 `json:"id"`

	// 字幕文件类型，目前暂只支持“SRT”。
	Type SubtitleType `json:"type"`

	// 字幕语言类型。  取值如下： - CN：表示中文字幕。 - EN：表示英文字幕。
	Language SubtitleLanguage `json:"language"`

	// 字幕文件的MD5值。
	Md5 *string `json:"md5,omitempty"`

	// 字幕描述。
	Description *string `json:"description,omitempty"`
}

func (o Subtitle) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Subtitle struct{}"
	}

	return strings.Join([]string{"Subtitle", string(data)}, " ")
}

type SubtitleType struct {
	value string
}

type SubtitleTypeEnum struct {
	SRT SubtitleType
}

func GetSubtitleTypeEnum() SubtitleTypeEnum {
	return SubtitleTypeEnum{
		SRT: SubtitleType{
			value: "SRT",
		},
	}
}

func (c SubtitleType) Value() string {
	return c.value
}

func (c SubtitleType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SubtitleType) UnmarshalJSON(b []byte) error {
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

type SubtitleLanguage struct {
	value string
}

type SubtitleLanguageEnum struct {
	CN SubtitleLanguage
	EN SubtitleLanguage
}

func GetSubtitleLanguageEnum() SubtitleLanguageEnum {
	return SubtitleLanguageEnum{
		CN: SubtitleLanguage{
			value: "CN",
		},
		EN: SubtitleLanguage{
			value: "EN",
		},
	}
}

func (c SubtitleLanguage) Value() string {
	return c.value
}

func (c SubtitleLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SubtitleLanguage) UnmarshalJSON(b []byte) error {
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
