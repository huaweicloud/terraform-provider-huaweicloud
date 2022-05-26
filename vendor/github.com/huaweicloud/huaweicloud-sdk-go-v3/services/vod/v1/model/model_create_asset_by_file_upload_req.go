package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CreateAssetByFileUploadReq struct {

	// 媒资标题，长度不超过128个字节，UTF-8编码。
	Title string `json:"title"`

	// 视频描述，长度不超过1024个字节。
	Description *string `json:"description,omitempty"`

	// 音视频文件名，长度不超过128个字节。  文件名后缀可选。
	VideoName string `json:"video_name"`

	// 上传音视频文件的格式。  取值如下： - 视频文件：MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG、3GP、ASF、MKV、HLS - 音频文件：MP3、OGG、WAV、WMA、APE、FLAC、AAC、AC3、MMF、AMR、M4A、M4R、WV、MP2  若上传格式为音频文件，则不支持转码、添加水印和字幕。
	VideoType string `json:"video_type"`

	// 媒资分类ID。  您可以调用[创建媒资分类](https://support.huaweicloud.com/api-vod/vod_04_0028.html)接口或在点播控制台的[分类设置](https://support.huaweicloud.com/usermanual-vod/vod010006.html)中创建对应的媒资分类，并获取分类ID。  > 若不设置或者设置为-1，则上传的音视频归类到系统预置的“其它”分类中。
	CategoryId *string `json:"category_id,omitempty"`

	// 视频文件MD5值。  建议参考[媒资上传和更新](https://support.huaweicloud.com/api-vod/vod_04_0212.html)生成对应的MD5值。
	VideoMd5 *string `json:"video_md5,omitempty"`

	// 封面图片文件类型。  取值如下： - JPG - PNG  上传后的封面名称是固定的，后缀名为封面类型缩写。例如cover0.jpg，cover1.png 若不指定类型，则封面文件无后缀名。  > 如果设置了图片格式，则不会执行首帧截图作为封面动作，需自行上传封面。
	CoverType *CreateAssetByFileUploadReqCoverType `json:"cover_type,omitempty"`

	// 封面文件MD5值
	CoverMd5 *string `json:"cover_md5,omitempty"`

	// 字幕文件信息
	Subtitles *[]Subtitle `json:"subtitles,omitempty"`

	// 视频标签。  单个标签不超过16个字节，最多不超过16个标签。  多个用逗号分隔，UTF8编码。
	Tags *string `json:"tags,omitempty"`

	// 是否自动发布。  取值如下： - 0：表示不自动发布。 - 1：表示自动发布。  默认值：0。
	AutoPublish *CreateAssetByFileUploadReqAutoPublish `json:"auto_publish,omitempty"`

	// 转码模板组名称。  若不为空，则使用指定的转码模板对上传的音视频进行转码，您可以在视频点播控制台配置转码模板，具体请参见[转码设置](https://support.huaweicloud.com/usermanual-vod/vod_01_0072.html)。  > 若同时设置了“**template_group_name**”和“**workflow_name**”字段，则“**template_group_name**”字段生效。
	TemplateGroupName *string `json:"template_group_name,omitempty"`

	// 是否自动加密。  取值如下： - 0：表示不加密。 - 1：表示需要加密。  默认值：0。  加密与转码必须要一起进行，当需要加密时，转码参数不能为空，且转码输出格式必须要为HLS。
	AutoEncrypt *int32 `json:"auto_encrypt,omitempty"`

	// 是否自动预热到CDN。  取值如下： - 0：表示不自动预热。 - 1：表示自动预热。  默认值：0。
	AutoPreheat *string `json:"auto_preheat,omitempty"`

	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`

	Review *Review `json:"review,omitempty"`

	// 工作流名称。  若不为空，则使用指定的工作流对上传的音视频进行处理，您可以在视频点播控制台配置工作流，具体请参见[工作流设置](https://support.huaweicloud.com/usermanual-vod/vod010041.html)。
	WorkflowName *string `json:"workflow_name,omitempty"`
}

func (o CreateAssetByFileUploadReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAssetByFileUploadReq struct{}"
	}

	return strings.Join([]string{"CreateAssetByFileUploadReq", string(data)}, " ")
}

type CreateAssetByFileUploadReqCoverType struct {
	value string
}

type CreateAssetByFileUploadReqCoverTypeEnum struct {
	JPG CreateAssetByFileUploadReqCoverType
	PNG CreateAssetByFileUploadReqCoverType
}

func GetCreateAssetByFileUploadReqCoverTypeEnum() CreateAssetByFileUploadReqCoverTypeEnum {
	return CreateAssetByFileUploadReqCoverTypeEnum{
		JPG: CreateAssetByFileUploadReqCoverType{
			value: "JPG",
		},
		PNG: CreateAssetByFileUploadReqCoverType{
			value: "PNG",
		},
	}
}

func (c CreateAssetByFileUploadReqCoverType) Value() string {
	return c.value
}

func (c CreateAssetByFileUploadReqCoverType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateAssetByFileUploadReqCoverType) UnmarshalJSON(b []byte) error {
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

type CreateAssetByFileUploadReqAutoPublish struct {
	value int32
}

type CreateAssetByFileUploadReqAutoPublishEnum struct {
	E_0 CreateAssetByFileUploadReqAutoPublish
	E_1 CreateAssetByFileUploadReqAutoPublish
}

func GetCreateAssetByFileUploadReqAutoPublishEnum() CreateAssetByFileUploadReqAutoPublishEnum {
	return CreateAssetByFileUploadReqAutoPublishEnum{
		E_0: CreateAssetByFileUploadReqAutoPublish{
			value: 0,
		}, E_1: CreateAssetByFileUploadReqAutoPublish{
			value: 1,
		},
	}
}

func (c CreateAssetByFileUploadReqAutoPublish) Value() int32 {
	return c.value
}

func (c CreateAssetByFileUploadReqAutoPublish) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateAssetByFileUploadReqAutoPublish) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(int32)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to int32 error")
	}
}
