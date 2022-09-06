package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type VideoInfo struct {

	// 画质<br/> 4K默认分辨率3840*2160，码率8000kbit/s<br/> 2K默认分辨率2560*1440，码率7000kbit/s<br/> FULL_HD默认分辨率1920*1080，码率3000kbit/s<br/> HD默认分辨率1280*720，码率1000kbit/s<br/> SD默认分辨率854*480，码率600kbit/s<br/> FLUENT默认分辨率480*270，码率300kbit/s<br/>
	Quality VideoInfoQuality `json:"quality"`

	// 视频宽度<br/>
	Width *int32 `json:"width,omitempty"`

	// 视频高度<br/>
	Height *int32 `json:"height,omitempty"`

	// 码率,单位：kbit/s<br/>
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 帧率（默认为0，0代表自适应，单位是帧每秒）<br/>
	FrameRate *int32 `json:"frame_rate,omitempty"`
}

func (o VideoInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoInfo struct{}"
	}

	return strings.Join([]string{"VideoInfo", string(data)}, " ")
}

type VideoInfoQuality struct {
	value string
}

type VideoInfoQualityEnum struct {
	FULL_HD VideoInfoQuality
	HD      VideoInfoQuality
	SD      VideoInfoQuality
	FLUENT  VideoInfoQuality
	E_2_K   VideoInfoQuality
	E_4_K   VideoInfoQuality
	UNKNOW  VideoInfoQuality
}

func GetVideoInfoQualityEnum() VideoInfoQualityEnum {
	return VideoInfoQualityEnum{
		FULL_HD: VideoInfoQuality{
			value: "FULL_HD",
		},
		HD: VideoInfoQuality{
			value: "HD",
		},
		SD: VideoInfoQuality{
			value: "SD",
		},
		FLUENT: VideoInfoQuality{
			value: "FLUENT",
		},
		E_2_K: VideoInfoQuality{
			value: "2K",
		},
		E_4_K: VideoInfoQuality{
			value: "4K",
		},
		UNKNOW: VideoInfoQuality{
			value: "UNKNOW",
		},
	}
}

func (c VideoInfoQuality) Value() string {
	return c.value
}

func (c VideoInfoQuality) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *VideoInfoQuality) UnmarshalJSON(b []byte) error {
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
