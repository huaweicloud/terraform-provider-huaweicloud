package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type VideoParameters struct {

	// 输出策略。  取值如下： - discard - transcode  >- 当视频参数中的“output_policy”为\"discard\"，且音频参数中的“output_policy”为“transcode”时，表示只输出音频。 >- 当视频参数中的“output_policy”为\"transcode\"，且音频参数中的“output_policy”为“discard”时，表示只输出视频。 >- 同时为\"discard\"时不合法。 >- 同时为“transcode”时，表示输出音视频。
	OutputPolicy *VideoParametersOutputPolicy `json:"output_policy,omitempty"`

	// 视频编码格式。  取值如下： - 1：VIDEO_CODEC_H264 - 2：VIDEO_CODEC_H265
	Codec *int32 `json:"codec,omitempty"`

	// 输出平均码率。  取值范围：0或[40,30000]之间的整数。  单位：kbit/s  若设置为0，则输出平均码率为自适应值。
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 编码档次  取值如下： - 1：VIDEO_PROFILE_H264_BASE - 2：VIDEO_PROFILE_H264_MAIN - 3：VIDEO_PROFILE_H264_HIGH - 4：VIDEO_PROFILE_H265_MAIN
	Profile *int32 `json:"profile,omitempty"`

	// 编码级别  取值如下： - 1：VIDEO_LEVEL_1_0 - 2：VIDEO_LEVEL_1_1 - 3：VIDEO_LEVEL_1_2 - 4：VIDEO_LEVEL_1_3 - 5：VIDEO_LEVEL_2_0 - 6：VIDEO_LEVEL_2_1 - 7：VIDEO_LEVEL_2_2 - 8：VIDEO_LEVEL_3_0 - 9：VIDEO_LEVEL_3_1 - 10：VIDEO_LEVEL_3_2 - 11：VIDEO_LEVEL_4_0 - 12：VIDEO_LEVEL_4_1 - 13：VIDEO_LEVEL_4_2 - 14：VIDEO_LEVEL_5_0 - 15：VIDEO_LEVEL_5_1 - 16：VIDEO_LEVEL_x_x
	Level *int32 `json:"level,omitempty"`

	// 编码质量等级  取值如下： - 1：VIDEO_PRESET_HSPEED2 (只用于h.265, h.265 default) - 2：VIDEO_PRESET_HSPEED (只用于h.265) - 3：VIDEO_PRESET_NORMAL (h264/h.265可用，h.264 default)
	Preset *int32 `json:"preset,omitempty"`

	// 最大参考帧数。  取值范围： - H264：[1，8]，默认值为4 。 - H265：固定值4。  单位：帧。
	RefFramesCount *int32 `json:"ref_frames_count,omitempty"`

	// I帧最大间隔  取值范围：[2，10]。  默认值：5。  单位：秒。
	MaxIframesInterval *int32 `json:"max_iframes_interval,omitempty"`

	// 最大B帧间隔。  取值范围： - H264：[0，7]，默认值为4。 - H265：[0，7]，默认值为7。  单位：帧。
	BframesCount *int32 `json:"bframes_count,omitempty"`

	// 帧率。  取值范围：0或[5,60]，0表示自适应。  单位：帧每秒。  > 若设置的帧率不在取值范围内，则自动调整为0，若设置的帧率高于片源帧率，则自动调整为片源帧率。
	FrameRate *int32 `json:"frame_rate,omitempty"`

	// 视频宽度（单位：像素）  - H264：范围[32,4096]，必须为2的倍数 - H265：范围[320,4096]，必须是4的倍数
	Width *int32 `json:"width,omitempty"`

	// 视频高度（单位：像素）  - H264：范围[32,2880]，必须为2的倍数 - H265：范围[240,2880] ，必须是4的倍数
	Height *int32 `json:"height,omitempty"`

	// 黑边剪裁类型  - 0：不开启黑边剪裁 - 1：开启黑边剪裁，低复杂度算法，针对长视频（>5分钟） - 2：开启黑边剪裁，高复杂度算法，针对短视频（<=5分钟）
	BlackCut *int32 `json:"black_cut,omitempty"`
}

func (o VideoParameters) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoParameters struct{}"
	}

	return strings.Join([]string{"VideoParameters", string(data)}, " ")
}

type VideoParametersOutputPolicy struct {
	value string
}

type VideoParametersOutputPolicyEnum struct {
	TRANSCODE VideoParametersOutputPolicy
	DISCARD   VideoParametersOutputPolicy
	COPY      VideoParametersOutputPolicy
}

func GetVideoParametersOutputPolicyEnum() VideoParametersOutputPolicyEnum {
	return VideoParametersOutputPolicyEnum{
		TRANSCODE: VideoParametersOutputPolicy{
			value: "transcode",
		},
		DISCARD: VideoParametersOutputPolicy{
			value: "discard",
		},
		COPY: VideoParametersOutputPolicy{
			value: "copy",
		},
	}
}

func (c VideoParametersOutputPolicy) Value() string {
	return c.value
}

func (c VideoParametersOutputPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *VideoParametersOutputPolicy) UnmarshalJSON(b []byte) error {
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
