package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// MetaData 视频的元数据信息。  经过视频解析后产生，包括封装格式、大小、分辨率、码率、帧率。
type MetaData struct {

	// 视频封装格式。  取值如下： - MP4 - TS - MOV - MXF - MPG - FLV - WMV - MP3 - WMA - APE - FLAC - AAC - AC3 - MMF - AMR - M4A - M4R - OGG - WAV - WV - MP2 - AVI - F4V - M4V - MPEG - HLS - DASH
	PackType *MetaDataPackType `json:"pack_type,omitempty"`

	// 视频编码格式。  取值如下： - MPEG-2 - MPEG-4 - H.264 - H.265 - WMV - Vorbis - AAC - AC-3 - AMR - APE - FLAC - MP3 - MP2 - WMA - PCM - ADPCM - WavPack
	Codec *MetaDataCodec `json:"codec,omitempty"`

	// 视频时长。  若视频的原时长为非整数，则该字段值为原时长的向下取整。 若视频的原时长小于1，则该字段值为1。
	Duration *int64 `json:"duration,omitempty"`

	// 视频文件大小。  单位：字节。
	VideoSize *int64 `json:"video_size,omitempty"`

	// 视频宽度（单位：像素）。 - 编码为H.264的取值范围：[32,3840]之间2的倍数。 - 编码为H.265的取值范围：[320,3840]之间4的倍数。
	Width *int64 `json:"width,omitempty"`

	// 视频高度（单位：像素）。 - 编码为H.264的取值范围：[32,2160]之间2的倍数 。 - 编码为H.265的取值范围：[240,2160]之间4的倍数。
	Hight *int64 `json:"hight,omitempty"`

	// 视频平均码率。
	BitRate *int64 `json:"bit_rate,omitempty"`

	// 帧率（单位：帧每秒）。  取值如下： - FRAMERATE_AUTO = 1, - FRAMERATE_10 = 2, - FRAMERATE_15 = 3, - FRAMERATE_2397 = 4, // 23.97 fps - FRAMERATE_24 = 5, - FRAMERATE_25 = 6, - FRAMERATE_2997 = 7, // 29.97 fps - FRAMERATE_30 = 8, - FRAMERATE_50 = 9, - FRAMERATE_60 = 10  默认值：1。  单位：帧每秒。
	FrameRate *int64 `json:"frame_rate,omitempty"`

	// 清晰度。  取值如下： - FULL_HD：超高清 - HD：高清 - SD：标清 - FLUENT：流畅 - AD：自适应 - 2K - 4K
	Quality *string `json:"quality,omitempty"`

	// 音频的声道数。
	AudioChannels *int32 `json:"audio_channels,omitempty"`
}

func (o MetaData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetaData struct{}"
	}

	return strings.Join([]string{"MetaData", string(data)}, " ")
}

type MetaDataPackType struct {
	value string
}

type MetaDataPackTypeEnum struct {
	MP4  MetaDataPackType
	TS   MetaDataPackType
	MOV  MetaDataPackType
	MXF  MetaDataPackType
	MPG  MetaDataPackType
	FLV  MetaDataPackType
	WMV  MetaDataPackType
	MP3  MetaDataPackType
	WMA  MetaDataPackType
	APE  MetaDataPackType
	FLAC MetaDataPackType
	AAC  MetaDataPackType
	AC3  MetaDataPackType
	MMF  MetaDataPackType
	AMR  MetaDataPackType
	M4_A MetaDataPackType
	M4_R MetaDataPackType
	OGG  MetaDataPackType
	WAV  MetaDataPackType
	WV   MetaDataPackType
	MP2  MetaDataPackType
	AVI  MetaDataPackType
	F4_V MetaDataPackType
	M4_V MetaDataPackType
	MPEG MetaDataPackType
	HLS  MetaDataPackType
	DASH MetaDataPackType
}

func GetMetaDataPackTypeEnum() MetaDataPackTypeEnum {
	return MetaDataPackTypeEnum{
		MP4: MetaDataPackType{
			value: "MP4",
		},
		TS: MetaDataPackType{
			value: "TS",
		},
		MOV: MetaDataPackType{
			value: "MOV",
		},
		MXF: MetaDataPackType{
			value: "MXF",
		},
		MPG: MetaDataPackType{
			value: "MPG",
		},
		FLV: MetaDataPackType{
			value: "FLV",
		},
		WMV: MetaDataPackType{
			value: "WMV",
		},
		MP3: MetaDataPackType{
			value: "MP3",
		},
		WMA: MetaDataPackType{
			value: "WMA",
		},
		APE: MetaDataPackType{
			value: "APE",
		},
		FLAC: MetaDataPackType{
			value: "FLAC",
		},
		AAC: MetaDataPackType{
			value: "AAC",
		},
		AC3: MetaDataPackType{
			value: "AC3",
		},
		MMF: MetaDataPackType{
			value: "MMF",
		},
		AMR: MetaDataPackType{
			value: "AMR",
		},
		M4_A: MetaDataPackType{
			value: "M4A",
		},
		M4_R: MetaDataPackType{
			value: "M4R",
		},
		OGG: MetaDataPackType{
			value: "OGG",
		},
		WAV: MetaDataPackType{
			value: "WAV",
		},
		WV: MetaDataPackType{
			value: "WV",
		},
		MP2: MetaDataPackType{
			value: "MP2",
		},
		AVI: MetaDataPackType{
			value: "AVI",
		},
		F4_V: MetaDataPackType{
			value: "F4V",
		},
		M4_V: MetaDataPackType{
			value: "M4V",
		},
		MPEG: MetaDataPackType{
			value: "MPEG",
		},
		HLS: MetaDataPackType{
			value: "HLS",
		},
		DASH: MetaDataPackType{
			value: "DASH",
		},
	}
}

func (c MetaDataPackType) Value() string {
	return c.value
}

func (c MetaDataPackType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MetaDataPackType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type MetaDataCodec struct {
	value string
}

type MetaDataCodecEnum struct {
	MPEG_2   MetaDataCodec
	MPEG_4   MetaDataCodec
	H_264    MetaDataCodec
	H_265    MetaDataCodec
	WMV      MetaDataCodec
	VORBIS   MetaDataCodec
	AAC      MetaDataCodec
	EAC_3    MetaDataCodec
	AC_3     MetaDataCodec
	AMR      MetaDataCodec
	APE      MetaDataCodec
	FLAC     MetaDataCodec
	MP3      MetaDataCodec
	MP2      MetaDataCodec
	WMA      MetaDataCodec
	PCM      MetaDataCodec
	ADPCM    MetaDataCodec
	WAV_PACK MetaDataCodec
	HEAAC    MetaDataCodec
	UNKNOWN  MetaDataCodec
}

func GetMetaDataCodecEnum() MetaDataCodecEnum {
	return MetaDataCodecEnum{
		MPEG_2: MetaDataCodec{
			value: "MPEG-2",
		},
		MPEG_4: MetaDataCodec{
			value: "MPEG-4",
		},
		H_264: MetaDataCodec{
			value: "H.264",
		},
		H_265: MetaDataCodec{
			value: "H.265",
		},
		WMV: MetaDataCodec{
			value: "WMV",
		},
		VORBIS: MetaDataCodec{
			value: "Vorbis",
		},
		AAC: MetaDataCodec{
			value: "AAC",
		},
		EAC_3: MetaDataCodec{
			value: "EAC-3",
		},
		AC_3: MetaDataCodec{
			value: "AC-3",
		},
		AMR: MetaDataCodec{
			value: "AMR",
		},
		APE: MetaDataCodec{
			value: "APE",
		},
		FLAC: MetaDataCodec{
			value: "FLAC",
		},
		MP3: MetaDataCodec{
			value: "MP3",
		},
		MP2: MetaDataCodec{
			value: "MP2",
		},
		WMA: MetaDataCodec{
			value: "WMA",
		},
		PCM: MetaDataCodec{
			value: "PCM",
		},
		ADPCM: MetaDataCodec{
			value: "ADPCM",
		},
		WAV_PACK: MetaDataCodec{
			value: "WavPack",
		},
		HEAAC: MetaDataCodec{
			value: "HEAAC",
		},
		UNKNOWN: MetaDataCodec{
			value: "UNKNOWN",
		},
	}
}

func (c MetaDataCodec) Value() string {
	return c.value
}

func (c MetaDataCodec) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MetaDataCodec) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
