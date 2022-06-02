package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type UploadAssetReq struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	// 视频文件MD5值。  建议参考[媒资上传和更新](https://support.huaweicloud.com/api-vod/vod_04_0212.html)生成对应的MD5值。
	VideoMd5 *string `json:"video_md5,omitempty"`

	// 视频文件名。  文件名后缀为可选。
	VideoName *string `json:"video_name,omitempty"`

	// 视频文件类型。 取值为MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG、3GP、ASF、MKV
	VideoType *UploadAssetReqVideoType `json:"video_type,omitempty"`

	// 封面ID。  取值范围：[0,7]。  当前只支持一张封面，只能设置为0。
	CoverId *int32 `json:"cover_id,omitempty"`

	// 封面图片格式类型。  取值如下： - JPG - PNG
	CoverType *UploadAssetReqCoverType `json:"cover_type,omitempty"`

	// 封面文件的MD5值。
	CoverMd5 *string `json:"cover_md5,omitempty"`

	// 字幕文件信息
	Subtitles *[]Subtitle `json:"subtitles,omitempty"`
}

func (o UploadAssetReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadAssetReq struct{}"
	}

	return strings.Join([]string{"UploadAssetReq", string(data)}, " ")
}

type UploadAssetReqVideoType struct {
	value string
}

type UploadAssetReqVideoTypeEnum struct {
	MP4    UploadAssetReqVideoType
	TS     UploadAssetReqVideoType
	MOV    UploadAssetReqVideoType
	MXF    UploadAssetReqVideoType
	MPG    UploadAssetReqVideoType
	FLV    UploadAssetReqVideoType
	WMV    UploadAssetReqVideoType
	HLS    UploadAssetReqVideoType
	MP3    UploadAssetReqVideoType
	WMA    UploadAssetReqVideoType
	APE    UploadAssetReqVideoType
	FLAC   UploadAssetReqVideoType
	AAC    UploadAssetReqVideoType
	AC3    UploadAssetReqVideoType
	MMF    UploadAssetReqVideoType
	AMR    UploadAssetReqVideoType
	M4_A   UploadAssetReqVideoType
	M4_R   UploadAssetReqVideoType
	OGG    UploadAssetReqVideoType
	WAV    UploadAssetReqVideoType
	WV     UploadAssetReqVideoType
	MP2    UploadAssetReqVideoType
	AVI    UploadAssetReqVideoType
	F4_V   UploadAssetReqVideoType
	M4_V   UploadAssetReqVideoType
	MPEG   UploadAssetReqVideoType
	UNKNOW UploadAssetReqVideoType
}

func GetUploadAssetReqVideoTypeEnum() UploadAssetReqVideoTypeEnum {
	return UploadAssetReqVideoTypeEnum{
		MP4: UploadAssetReqVideoType{
			value: "MP4",
		},
		TS: UploadAssetReqVideoType{
			value: "TS",
		},
		MOV: UploadAssetReqVideoType{
			value: "MOV",
		},
		MXF: UploadAssetReqVideoType{
			value: "MXF",
		},
		MPG: UploadAssetReqVideoType{
			value: "MPG",
		},
		FLV: UploadAssetReqVideoType{
			value: "FLV",
		},
		WMV: UploadAssetReqVideoType{
			value: "WMV",
		},
		HLS: UploadAssetReqVideoType{
			value: "HLS",
		},
		MP3: UploadAssetReqVideoType{
			value: "MP3",
		},
		WMA: UploadAssetReqVideoType{
			value: "WMA",
		},
		APE: UploadAssetReqVideoType{
			value: "APE",
		},
		FLAC: UploadAssetReqVideoType{
			value: "FLAC",
		},
		AAC: UploadAssetReqVideoType{
			value: "AAC",
		},
		AC3: UploadAssetReqVideoType{
			value: "AC3",
		},
		MMF: UploadAssetReqVideoType{
			value: "MMF",
		},
		AMR: UploadAssetReqVideoType{
			value: "AMR",
		},
		M4_A: UploadAssetReqVideoType{
			value: "M4A",
		},
		M4_R: UploadAssetReqVideoType{
			value: "M4R",
		},
		OGG: UploadAssetReqVideoType{
			value: "OGG",
		},
		WAV: UploadAssetReqVideoType{
			value: "WAV",
		},
		WV: UploadAssetReqVideoType{
			value: "WV",
		},
		MP2: UploadAssetReqVideoType{
			value: "MP2",
		},
		AVI: UploadAssetReqVideoType{
			value: "AVI",
		},
		F4_V: UploadAssetReqVideoType{
			value: "F4V",
		},
		M4_V: UploadAssetReqVideoType{
			value: "M4V",
		},
		MPEG: UploadAssetReqVideoType{
			value: "MPEG",
		},
		UNKNOW: UploadAssetReqVideoType{
			value: "UNKNOW",
		},
	}
}

func (c UploadAssetReqVideoType) Value() string {
	return c.value
}

func (c UploadAssetReqVideoType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UploadAssetReqVideoType) UnmarshalJSON(b []byte) error {
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

type UploadAssetReqCoverType struct {
	value string
}

type UploadAssetReqCoverTypeEnum struct {
	JPG UploadAssetReqCoverType
	PNG UploadAssetReqCoverType
}

func GetUploadAssetReqCoverTypeEnum() UploadAssetReqCoverTypeEnum {
	return UploadAssetReqCoverTypeEnum{
		JPG: UploadAssetReqCoverType{
			value: "JPG",
		},
		PNG: UploadAssetReqCoverType{
			value: "PNG",
		},
	}
}

func (c UploadAssetReqCoverType) Value() string {
	return c.value
}

func (c UploadAssetReqCoverType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UploadAssetReqCoverType) UnmarshalJSON(b []byte) error {
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
