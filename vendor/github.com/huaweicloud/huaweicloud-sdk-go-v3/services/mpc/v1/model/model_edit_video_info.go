package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type EditVideoInfo struct {

	// 剪辑输出视频参数的参照物。取值如下： - MAX，以输入片源中最大分辨率的视频参数作为输出参照。 - MIN，以输入片源中最小分辨率的视频参数作为输出参照。 - CUSTOM，自定义视频输出参数，使用该参数时，所有视频参数必填。- SHORT_HEIGHT_SHORT_WIDTH，当edit_type为MIX时，只能使用该值。
	Reference *EditVideoInfoReference `json:"reference,omitempty"`

	// 视频宽度。
	Width *int32 `json:"width,omitempty"`

	// 视频高度。
	Height *int32 `json:"height,omitempty"`

	// 视频频编码格式。
	Codec *EditVideoInfoCodec `json:"codec,omitempty"`

	// 视频码率，单位: bit/s
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 帧率。
	FrameRate *int32 `json:"frame_rate,omitempty"`
}

func (o EditVideoInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EditVideoInfo struct{}"
	}

	return strings.Join([]string{"EditVideoInfo", string(data)}, " ")
}

type EditVideoInfoReference struct {
	value string
}

type EditVideoInfoReferenceEnum struct {
	MAX                      EditVideoInfoReference
	MIN                      EditVideoInfoReference
	CUSTOM                   EditVideoInfoReference
	SHORT_HEIGHT_SHORT_WIDTH EditVideoInfoReference
}

func GetEditVideoInfoReferenceEnum() EditVideoInfoReferenceEnum {
	return EditVideoInfoReferenceEnum{
		MAX: EditVideoInfoReference{
			value: "MAX",
		},
		MIN: EditVideoInfoReference{
			value: "MIN",
		},
		CUSTOM: EditVideoInfoReference{
			value: "CUSTOM",
		},
		SHORT_HEIGHT_SHORT_WIDTH: EditVideoInfoReference{
			value: "SHORT_HEIGHT_SHORT_WIDTH",
		},
	}
}

func (c EditVideoInfoReference) Value() string {
	return c.value
}

func (c EditVideoInfoReference) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EditVideoInfoReference) UnmarshalJSON(b []byte) error {
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

type EditVideoInfoCodec struct {
	value string
}

type EditVideoInfoCodecEnum struct {
	H264 EditVideoInfoCodec
	H265 EditVideoInfoCodec
}

func GetEditVideoInfoCodecEnum() EditVideoInfoCodecEnum {
	return EditVideoInfoCodecEnum{
		H264: EditVideoInfoCodec{
			value: "H264",
		},
		H265: EditVideoInfoCodec{
			value: "H265",
		},
	}
}

func (c EditVideoInfoCodec) Value() string {
	return c.value
}

func (c EditVideoInfoCodec) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EditVideoInfoCodec) UnmarshalJSON(b []byte) error {
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
