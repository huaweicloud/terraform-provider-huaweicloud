package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type VideoTypeRef struct {

	// 转存的音视频文件类型。  取值如下： - 视频文件：MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG、3GP、ASF、MKV、HLS - 音频文件：MP3、OGG、WAV、WMA、APE、FLAC、AAC、AC3、MMF、AMR、M4A、M4R、WV、MP2  若上传格式为音频文件，则不支持转码、添加水印和字幕。  > 当**video_type**选择HLS时，**storage_mode**（存储模式）需选择存储在租户桶，且输出路径设置为和输入路径在同一个目录。
	VideoType VideoTypeRefVideoType `json:"video_type"`

	// 媒资标题，长度不超过128个字节，UTF-8编码。
	Title string `json:"title"`

	// 视频描述，长度不超过1024个字节。
	Description *string `json:"description,omitempty"`

	// 媒资分类ID。  您可以调用[创建媒资分类](https://support.huaweicloud.com/api-vod/vod_04_0028.html)接口或在点播控制台的[分类设置](https://support.huaweicloud.com/usermanual-vod/vod010006.html)中创建对应的媒资分类，并获取分类ID。  > 若不设置或者设置为-1，则上传的音视频归类到系统预置的“其它”分类中。
	CategoryId *int32 `json:"category_id,omitempty"`

	// 视频标签。  单个标签不超过16个字节，最多不超过16个标签。  多个用逗号分隔，UTF8编码。
	Tags *string `json:"tags,omitempty"`

	// 是否自动发布。  取值如下： - 0：表示不自动发布。 - 1：表示自动发布。  默认值：0。
	AutoPublish *int32 `json:"auto_publish,omitempty"`

	// 转码模板组名称。  若不为空，则使用指定的转码模板对上传的音视频进行转码，您可以在视频点播控制台配置转码模板，具体请参见[转码设置](https://support.huaweicloud.com/usermanual-vod/vod_01_0072.html)。  > 若同时设置了“**template_group_name**”和“**workflow_name**”字段，则“**template_group_name**”字段生效。
	TemplateGroupName *string `json:"template_group_name,omitempty"`

	// 是否自动加密。  取值如下： - 0：表示不加密。 - 1：表示需要加密。  默认值：0。  若设置为需要加密，则必须配置转码模板，且转码的输出格式是HLS。
	AutoEncrypt *int32 `json:"auto_encrypt,omitempty"`

	// 是否自动预热到CDN。  取值如下： - 0：表示不自动预热。 - 1：表示自动预热。  默认值：0。
	AutoPreheat *int32 `json:"auto_preheat,omitempty"`

	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`

	Review *Review `json:"review,omitempty"`

	// 工作流名称。  若不为空，则使用指定的工作流对上传的音视频进行处理，您可以在视频点播控制台配置工作流，具体请参见[工作流设置](https://support.huaweicloud.com/usermanual-vod/vod010041.html)。
	WorkflowName *string `json:"workflow_name,omitempty"`
}

func (o VideoTypeRef) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoTypeRef struct{}"
	}

	return strings.Join([]string{"VideoTypeRef", string(data)}, " ")
}

type VideoTypeRefVideoType struct {
	value string
}

type VideoTypeRefVideoTypeEnum struct {
	MP4    VideoTypeRefVideoType
	TS     VideoTypeRefVideoType
	MOV    VideoTypeRefVideoType
	MXF    VideoTypeRefVideoType
	MPG    VideoTypeRefVideoType
	FLV    VideoTypeRefVideoType
	WMV    VideoTypeRefVideoType
	AVI    VideoTypeRefVideoType
	M4_V   VideoTypeRefVideoType
	F4_V   VideoTypeRefVideoType
	MPEG   VideoTypeRefVideoType
	E_3_GP VideoTypeRefVideoType
	ASF    VideoTypeRefVideoType
	MKV    VideoTypeRefVideoType
	HLS    VideoTypeRefVideoType
	MP3    VideoTypeRefVideoType
	OGG    VideoTypeRefVideoType
	WAV    VideoTypeRefVideoType
	WMA    VideoTypeRefVideoType
	APE    VideoTypeRefVideoType
	FLAC   VideoTypeRefVideoType
	AAC    VideoTypeRefVideoType
	AC3    VideoTypeRefVideoType
	MMF    VideoTypeRefVideoType
	AMR    VideoTypeRefVideoType
	M4_A   VideoTypeRefVideoType
	M4_R   VideoTypeRefVideoType
	WV     VideoTypeRefVideoType
	MP2    VideoTypeRefVideoType
}

func GetVideoTypeRefVideoTypeEnum() VideoTypeRefVideoTypeEnum {
	return VideoTypeRefVideoTypeEnum{
		MP4: VideoTypeRefVideoType{
			value: "MP4",
		},
		TS: VideoTypeRefVideoType{
			value: "TS",
		},
		MOV: VideoTypeRefVideoType{
			value: "MOV",
		},
		MXF: VideoTypeRefVideoType{
			value: "MXF",
		},
		MPG: VideoTypeRefVideoType{
			value: "MPG",
		},
		FLV: VideoTypeRefVideoType{
			value: "FLV",
		},
		WMV: VideoTypeRefVideoType{
			value: "WMV",
		},
		AVI: VideoTypeRefVideoType{
			value: "AVI",
		},
		M4_V: VideoTypeRefVideoType{
			value: "M4V",
		},
		F4_V: VideoTypeRefVideoType{
			value: "F4V",
		},
		MPEG: VideoTypeRefVideoType{
			value: "MPEG",
		},
		E_3_GP: VideoTypeRefVideoType{
			value: "3GP",
		},
		ASF: VideoTypeRefVideoType{
			value: "ASF",
		},
		MKV: VideoTypeRefVideoType{
			value: "MKV",
		},
		HLS: VideoTypeRefVideoType{
			value: "HLS",
		},
		MP3: VideoTypeRefVideoType{
			value: "MP3",
		},
		OGG: VideoTypeRefVideoType{
			value: "OGG",
		},
		WAV: VideoTypeRefVideoType{
			value: "WAV",
		},
		WMA: VideoTypeRefVideoType{
			value: "WMA",
		},
		APE: VideoTypeRefVideoType{
			value: "APE",
		},
		FLAC: VideoTypeRefVideoType{
			value: "FLAC",
		},
		AAC: VideoTypeRefVideoType{
			value: "AAC",
		},
		AC3: VideoTypeRefVideoType{
			value: "AC3",
		},
		MMF: VideoTypeRefVideoType{
			value: "MMF",
		},
		AMR: VideoTypeRefVideoType{
			value: "AMR",
		},
		M4_A: VideoTypeRefVideoType{
			value: "M4A",
		},
		M4_R: VideoTypeRefVideoType{
			value: "M4R",
		},
		WV: VideoTypeRefVideoType{
			value: "WV",
		},
		MP2: VideoTypeRefVideoType{
			value: "MP2",
		},
	}
}

func (c VideoTypeRefVideoType) Value() string {
	return c.value
}

func (c VideoTypeRefVideoType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *VideoTypeRefVideoType) UnmarshalJSON(b []byte) error {
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
