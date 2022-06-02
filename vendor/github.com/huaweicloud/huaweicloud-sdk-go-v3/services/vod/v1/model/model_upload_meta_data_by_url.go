package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type UploadMetaDataByUrl struct {

	// 上传音视频文件的格式。  取值如下： - 视频文件：MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG、3GP、ASF、MKV、M3U8 - 音频文件：MP3、OGG、WAV、WMA、APE、FLAC、AAC、AC3、MMF、AMR、M4A、M4R、WV、MP2  若上传格式为音频文件，则不支持转码、添加水印和字幕。
	VideoType UploadMetaDataByUrlVideoType `json:"video_type"`

	// 媒资标题，长度不超过128个字节，UTF-8编码。
	Title string `json:"title"`

	// 音视频源文件URL。   > URL必须以扩展名结尾，暂只支持http和https协议。
	Url string `json:"url"`

	// 视频描述，长度不超过1024个字节。
	Description *string `json:"description,omitempty"`

	// 媒资分类ID。  您可以调用[创建媒资分类](https://support.huaweicloud.com/api-vod/vod_04_0028.html)接口或在点播控制台的[分类设置](https://support.huaweicloud.com/usermanual-vod/vod010006.html)中创建对应的媒资分类，并获取分类ID。  > 若不设置或者设置为-1，则上传的音视频归类到系统预置的“其它”分类中。
	CategoryId *int32 `json:"category_id,omitempty"`

	// 视频标签。  单个标签不超过16个字节，最多不超过16个标签。  多个用逗号分隔，UTF8编码。
	Tags *string `json:"tags,omitempty"`

	// 是否自动发布。  取值如下： - 0：表示不自动发布。 - 1：表示自动发布。  默认值：0。
	AutoPublish *int32 `json:"auto_publish,omitempty"`

	// 转码模板组名称。  若不为空，则使用指定的转码模板对上传的音视频进行转码，您可以在视频点播控制台配置转码模板，具体请参见[转码设置](https://support.huaweicloud.com/usermanual-vod/vod_01_0072.html)。  >若同时设置了“**template_group_name**”和“**workflow_name**”字段，则“**template_group_name**”字段生效。
	TemplateGroupName *string `json:"template_group_name,omitempty"`

	// 是否自动加密。  取值如下： - 0：表示不加密。 - 1：表示需要加密。  默认值：0。若设置为需要加密，则必须配置转码模板，且转码的输出格式是HLS。
	AutoEncrypt *int32 `json:"auto_encrypt,omitempty"`

	// 是否自动预热到CDN。  取值如下： - 0：表示不自动预热。 - 1：表示自动预热。  默认值：0。
	AutoPreheat *int32 `json:"auto_preheat,omitempty"`

	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`

	Review *Review `json:"review,omitempty"`

	// 工作流名称。  若不为空，则使用指定的工作流对上传的音视频进行处理，您可以在视频点播控制台配置工作流，具体请参见[工作流设置](https://support.huaweicloud.com/usermanual-vod/vod010041.html)。
	WorkflowName *string `json:"workflow_name,omitempty"`
}

func (o UploadMetaDataByUrl) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadMetaDataByUrl struct{}"
	}

	return strings.Join([]string{"UploadMetaDataByUrl", string(data)}, " ")
}

type UploadMetaDataByUrlVideoType struct {
	value string
}

type UploadMetaDataByUrlVideoTypeEnum struct {
	MP4    UploadMetaDataByUrlVideoType
	TS     UploadMetaDataByUrlVideoType
	MOV    UploadMetaDataByUrlVideoType
	MXF    UploadMetaDataByUrlVideoType
	MPG    UploadMetaDataByUrlVideoType
	FLV    UploadMetaDataByUrlVideoType
	WMV    UploadMetaDataByUrlVideoType
	AVI    UploadMetaDataByUrlVideoType
	M4_V   UploadMetaDataByUrlVideoType
	F4_V   UploadMetaDataByUrlVideoType
	MPEG   UploadMetaDataByUrlVideoType
	E_3_GP UploadMetaDataByUrlVideoType
	ASF    UploadMetaDataByUrlVideoType
	MKV    UploadMetaDataByUrlVideoType
	MP3    UploadMetaDataByUrlVideoType
	OGG    UploadMetaDataByUrlVideoType
	WAV    UploadMetaDataByUrlVideoType
	WMA    UploadMetaDataByUrlVideoType
	APE    UploadMetaDataByUrlVideoType
	FLAC   UploadMetaDataByUrlVideoType
	AAC    UploadMetaDataByUrlVideoType
	AC3    UploadMetaDataByUrlVideoType
	MMF    UploadMetaDataByUrlVideoType
	AMR    UploadMetaDataByUrlVideoType
	M4_A   UploadMetaDataByUrlVideoType
	M4_R   UploadMetaDataByUrlVideoType
	WV     UploadMetaDataByUrlVideoType
	MP2    UploadMetaDataByUrlVideoType
	M3_U8  UploadMetaDataByUrlVideoType
}

func GetUploadMetaDataByUrlVideoTypeEnum() UploadMetaDataByUrlVideoTypeEnum {
	return UploadMetaDataByUrlVideoTypeEnum{
		MP4: UploadMetaDataByUrlVideoType{
			value: "MP4",
		},
		TS: UploadMetaDataByUrlVideoType{
			value: "TS",
		},
		MOV: UploadMetaDataByUrlVideoType{
			value: "MOV",
		},
		MXF: UploadMetaDataByUrlVideoType{
			value: "MXF",
		},
		MPG: UploadMetaDataByUrlVideoType{
			value: "MPG",
		},
		FLV: UploadMetaDataByUrlVideoType{
			value: "FLV",
		},
		WMV: UploadMetaDataByUrlVideoType{
			value: "WMV",
		},
		AVI: UploadMetaDataByUrlVideoType{
			value: "AVI",
		},
		M4_V: UploadMetaDataByUrlVideoType{
			value: "M4V",
		},
		F4_V: UploadMetaDataByUrlVideoType{
			value: "F4V",
		},
		MPEG: UploadMetaDataByUrlVideoType{
			value: "MPEG",
		},
		E_3_GP: UploadMetaDataByUrlVideoType{
			value: "3GP",
		},
		ASF: UploadMetaDataByUrlVideoType{
			value: "ASF",
		},
		MKV: UploadMetaDataByUrlVideoType{
			value: "MKV",
		},
		MP3: UploadMetaDataByUrlVideoType{
			value: "MP3",
		},
		OGG: UploadMetaDataByUrlVideoType{
			value: "OGG",
		},
		WAV: UploadMetaDataByUrlVideoType{
			value: "WAV",
		},
		WMA: UploadMetaDataByUrlVideoType{
			value: "WMA",
		},
		APE: UploadMetaDataByUrlVideoType{
			value: "APE",
		},
		FLAC: UploadMetaDataByUrlVideoType{
			value: "FLAC",
		},
		AAC: UploadMetaDataByUrlVideoType{
			value: "AAC",
		},
		AC3: UploadMetaDataByUrlVideoType{
			value: "AC3",
		},
		MMF: UploadMetaDataByUrlVideoType{
			value: "MMF",
		},
		AMR: UploadMetaDataByUrlVideoType{
			value: "AMR",
		},
		M4_A: UploadMetaDataByUrlVideoType{
			value: "M4A",
		},
		M4_R: UploadMetaDataByUrlVideoType{
			value: "M4R",
		},
		WV: UploadMetaDataByUrlVideoType{
			value: "WV",
		},
		MP2: UploadMetaDataByUrlVideoType{
			value: "MP2",
		},
		M3_U8: UploadMetaDataByUrlVideoType{
			value: "M3U8",
		},
	}
}

func (c UploadMetaDataByUrlVideoType) Value() string {
	return c.value
}

func (c UploadMetaDataByUrlVideoType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UploadMetaDataByUrlVideoType) UnmarshalJSON(b []byte) error {
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
