package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

//
type PublishAssetFromObsReq struct {

	// 转存的音视频文件类型。  取值如下： - 视频文件：MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG、3GP、ASF、MKV、HLS - 音频文件：MP3、OGG、WAV、WMA、APE、FLAC、AAC、AC3、MMF、AMR、M4A、M4R、WV、MP2  若上传格式为音频文件，则不支持转码、添加水印和字幕。  > 当**video_type**选择HLS时，**storage_mode**（存储模式）需选择存储在租户桶，且输出路径设置为和输入路径在同一个目录。
	VideoType PublishAssetFromObsReqVideoType `json:"video_type"`

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

	Input *FileAddr `json:"input"`

	// 存储模式。  取值如下： - 0：表示视频拷贝到点播桶。 - 1：表示视频存储在租户桶。  默认值：0
	StorageMode *int32 `json:"storage_mode,omitempty"`

	// 输出桶名，“**storage_mode**”为1时必选。
	OutputBucket *string `json:"output_bucket,omitempty"`

	// 输出路径名，“**storage_mode**”为1时必选。
	OutputPath *string `json:"output_path,omitempty"`
}

func (o PublishAssetFromObsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PublishAssetFromObsReq struct{}"
	}

	return strings.Join([]string{"PublishAssetFromObsReq", string(data)}, " ")
}

type PublishAssetFromObsReqVideoType struct {
	value string
}

type PublishAssetFromObsReqVideoTypeEnum struct {
	MP4    PublishAssetFromObsReqVideoType
	TS     PublishAssetFromObsReqVideoType
	MOV    PublishAssetFromObsReqVideoType
	MXF    PublishAssetFromObsReqVideoType
	MPG    PublishAssetFromObsReqVideoType
	FLV    PublishAssetFromObsReqVideoType
	WMV    PublishAssetFromObsReqVideoType
	AVI    PublishAssetFromObsReqVideoType
	M4_V   PublishAssetFromObsReqVideoType
	F4_V   PublishAssetFromObsReqVideoType
	MPEG   PublishAssetFromObsReqVideoType
	E_3_GP PublishAssetFromObsReqVideoType
	ASF    PublishAssetFromObsReqVideoType
	MKV    PublishAssetFromObsReqVideoType
	HLS    PublishAssetFromObsReqVideoType
	MP3    PublishAssetFromObsReqVideoType
	OGG    PublishAssetFromObsReqVideoType
	WAV    PublishAssetFromObsReqVideoType
	WMA    PublishAssetFromObsReqVideoType
	APE    PublishAssetFromObsReqVideoType
	FLAC   PublishAssetFromObsReqVideoType
	AAC    PublishAssetFromObsReqVideoType
	AC3    PublishAssetFromObsReqVideoType
	MMF    PublishAssetFromObsReqVideoType
	AMR    PublishAssetFromObsReqVideoType
	M4_A   PublishAssetFromObsReqVideoType
	M4_R   PublishAssetFromObsReqVideoType
	WV     PublishAssetFromObsReqVideoType
	MP2    PublishAssetFromObsReqVideoType
}

func GetPublishAssetFromObsReqVideoTypeEnum() PublishAssetFromObsReqVideoTypeEnum {
	return PublishAssetFromObsReqVideoTypeEnum{
		MP4: PublishAssetFromObsReqVideoType{
			value: "MP4",
		},
		TS: PublishAssetFromObsReqVideoType{
			value: "TS",
		},
		MOV: PublishAssetFromObsReqVideoType{
			value: "MOV",
		},
		MXF: PublishAssetFromObsReqVideoType{
			value: "MXF",
		},
		MPG: PublishAssetFromObsReqVideoType{
			value: "MPG",
		},
		FLV: PublishAssetFromObsReqVideoType{
			value: "FLV",
		},
		WMV: PublishAssetFromObsReqVideoType{
			value: "WMV",
		},
		AVI: PublishAssetFromObsReqVideoType{
			value: "AVI",
		},
		M4_V: PublishAssetFromObsReqVideoType{
			value: "M4V",
		},
		F4_V: PublishAssetFromObsReqVideoType{
			value: "F4V",
		},
		MPEG: PublishAssetFromObsReqVideoType{
			value: "MPEG",
		},
		E_3_GP: PublishAssetFromObsReqVideoType{
			value: "3GP",
		},
		ASF: PublishAssetFromObsReqVideoType{
			value: "ASF",
		},
		MKV: PublishAssetFromObsReqVideoType{
			value: "MKV",
		},
		HLS: PublishAssetFromObsReqVideoType{
			value: "HLS",
		},
		MP3: PublishAssetFromObsReqVideoType{
			value: "MP3",
		},
		OGG: PublishAssetFromObsReqVideoType{
			value: "OGG",
		},
		WAV: PublishAssetFromObsReqVideoType{
			value: "WAV",
		},
		WMA: PublishAssetFromObsReqVideoType{
			value: "WMA",
		},
		APE: PublishAssetFromObsReqVideoType{
			value: "APE",
		},
		FLAC: PublishAssetFromObsReqVideoType{
			value: "FLAC",
		},
		AAC: PublishAssetFromObsReqVideoType{
			value: "AAC",
		},
		AC3: PublishAssetFromObsReqVideoType{
			value: "AC3",
		},
		MMF: PublishAssetFromObsReqVideoType{
			value: "MMF",
		},
		AMR: PublishAssetFromObsReqVideoType{
			value: "AMR",
		},
		M4_A: PublishAssetFromObsReqVideoType{
			value: "M4A",
		},
		M4_R: PublishAssetFromObsReqVideoType{
			value: "M4R",
		},
		WV: PublishAssetFromObsReqVideoType{
			value: "WV",
		},
		MP2: PublishAssetFromObsReqVideoType{
			value: "MP2",
		},
	}
}

func (c PublishAssetFromObsReqVideoType) Value() string {
	return c.value
}

func (c PublishAssetFromObsReqVideoType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PublishAssetFromObsReqVideoType) UnmarshalJSON(b []byte) error {
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
