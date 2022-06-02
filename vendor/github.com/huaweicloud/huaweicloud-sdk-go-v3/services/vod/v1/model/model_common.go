package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 模板信息。
type Common struct {

	// pvc开关。
	Pvc CommonPvc `json:"pvc"`

	// pvc版本。
	PvcVersion *string `json:"pvc_version,omitempty"`

	// 视频编码格式。
	VideoCodec *CommonVideoCodec `json:"video_codec,omitempty"`

	// 音频编码格式(有效值范围) - 1：AUDIO_CODECTYPE_AAC - 2：AUDIO_CODECTYPE_HEAAC1 - 3：AUDIO_CODECTYPE_HEAAC2 - 4：AUDIO_CODECTYPE_MP3  默认值为1。
	AudioCodec *CommonAudioCodec `json:"audio_codec,omitempty"`

	// 分片时长(默认为5秒)。
	HlsInterval *int32 `json:"hls_interval,omitempty"`
}

func (o Common) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Common struct{}"
	}

	return strings.Join([]string{"Common", string(data)}, " ")
}

type CommonPvc struct {
	value string
}

type CommonPvcEnum struct {
	E_0    CommonPvc
	E_1    CommonPvc
	E_2    CommonPvc
	UNKNOW CommonPvc
}

func GetCommonPvcEnum() CommonPvcEnum {
	return CommonPvcEnum{
		E_0: CommonPvc{
			value: "0",
		},
		E_1: CommonPvc{
			value: "1",
		},
		E_2: CommonPvc{
			value: "2",
		},
		UNKNOW: CommonPvc{
			value: "UNKNOW",
		},
	}
}

func (c CommonPvc) Value() string {
	return c.value
}

func (c CommonPvc) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CommonPvc) UnmarshalJSON(b []byte) error {
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

type CommonVideoCodec struct {
	value string
}

type CommonVideoCodecEnum struct {
	H264   CommonVideoCodec
	H265   CommonVideoCodec
	UNKNOW CommonVideoCodec
}

func GetCommonVideoCodecEnum() CommonVideoCodecEnum {
	return CommonVideoCodecEnum{
		H264: CommonVideoCodec{
			value: "H264",
		},
		H265: CommonVideoCodec{
			value: "H265",
		},
		UNKNOW: CommonVideoCodec{
			value: "UNKNOW",
		},
	}
}

func (c CommonVideoCodec) Value() string {
	return c.value
}

func (c CommonVideoCodec) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CommonVideoCodec) UnmarshalJSON(b []byte) error {
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

type CommonAudioCodec struct {
	value string
}

type CommonAudioCodecEnum struct {
	AAC    CommonAudioCodec
	HEAAC1 CommonAudioCodec
	HEAAC2 CommonAudioCodec
	MP3    CommonAudioCodec
}

func GetCommonAudioCodecEnum() CommonAudioCodecEnum {
	return CommonAudioCodecEnum{
		AAC: CommonAudioCodec{
			value: "AAC",
		},
		HEAAC1: CommonAudioCodec{
			value: "HEAAC1",
		},
		HEAAC2: CommonAudioCodec{
			value: "HEAAC2",
		},
		MP3: CommonAudioCodec{
			value: "MP3",
		},
	}
}

func (c CommonAudioCodec) Value() string {
	return c.value
}

func (c CommonAudioCodec) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CommonAudioCodec) UnmarshalJSON(b []byte) error {
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
