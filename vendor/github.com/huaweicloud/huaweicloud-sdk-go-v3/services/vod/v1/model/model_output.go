package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 转码输出数组。 - HLS或DASH格式：此数组的成员个数为n+1，n为转码输出路数。 - MP4格式：此数组的成员个数为n，n为转码输出路数。
type Output struct {

	// 协议类型。  取值如下： - hls - dash - mp4
	PlayType OutputPlayType `json:"play_type"`

	// 播放URL。
	Url string `json:"url"`

	// 标记流是否已被加密。  取值如下： - 0：表示未加密。 - 1：表示已被加密。
	Encrypted *int32 `json:"encrypted,omitempty"`

	// 清晰度。  取值如下： - FLUENT：流畅 - SD：标清 - HD：高清 - FULL_HD：超清
	Quality *OutputQuality `json:"quality,omitempty"`

	MetaData *MetaData `json:"meta_data"`
}

func (o Output) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Output struct{}"
	}

	return strings.Join([]string{"Output", string(data)}, " ")
}

type OutputPlayType struct {
	value string
}

type OutputPlayTypeEnum struct {
	HLS  OutputPlayType
	DASH OutputPlayType
	MP4  OutputPlayType
	MP3  OutputPlayType
	AAC  OutputPlayType
}

func GetOutputPlayTypeEnum() OutputPlayTypeEnum {
	return OutputPlayTypeEnum{
		HLS: OutputPlayType{
			value: "HLS",
		},
		DASH: OutputPlayType{
			value: "DASH",
		},
		MP4: OutputPlayType{
			value: "MP4",
		},
		MP3: OutputPlayType{
			value: "MP3",
		},
		AAC: OutputPlayType{
			value: "AAC",
		},
	}
}

func (c OutputPlayType) Value() string {
	return c.value
}

func (c OutputPlayType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OutputPlayType) UnmarshalJSON(b []byte) error {
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

type OutputQuality struct {
	value string
}

type OutputQualityEnum struct {
	FLUENT  OutputQuality
	SD      OutputQuality
	HD      OutputQuality
	FULL_HD OutputQuality
}

func GetOutputQualityEnum() OutputQualityEnum {
	return OutputQualityEnum{
		FLUENT: OutputQuality{
			value: "FLUENT",
		},
		SD: OutputQuality{
			value: "SD",
		},
		HD: OutputQuality{
			value: "HD",
		},
		FULL_HD: OutputQuality{
			value: "FULL_HD",
		},
	}
}

func (c OutputQuality) Value() string {
	return c.value
}

func (c OutputQuality) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OutputQuality) UnmarshalJSON(b []byte) error {
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
