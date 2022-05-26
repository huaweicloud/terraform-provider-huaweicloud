package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Video struct {

	// 输出策略。  取值如下： - discard - transcode  >- 当视频参数中的“output_policy”为\"discard\"，且音频参数中的“output_policy”为“transcode”时，表示只输出音频。 >- 当视频参数中的“output_policy”为\"transcode\"，且音频参数中的“output_policy”为“discard”时，表示只输出视频。 >- 同时为\"discard\"时不合法。 >- 同时为“transcode”时，表示输出音视频。
	OutputPolicy *VideoOutputPolicy `json:"output_policy,omitempty"`

	// 视频编码格式。  取值如下：  - 1：表示H.264。 - 2：表示H.265。
	Codec *int32 `json:"codec,omitempty"`

	// 输出平均码率。  取值范围：0或[40,30000]之间的整数。  单位：kbit/s  若设置为0，则输出平均码率为自适应值。
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 编码档次，建议设为3。  取值如下： - 1：VIDEO_PROFILE_H264_BASE - 2：VIDEO_PROFILE_H264_MAIN - 3：VIDEO_PROFILE_H264_HIGH - 4：VIDEO_PROFILE_H265_MAIN
	Profile *int32 `json:"profile,omitempty"`

	// 编码级别。  取值如下： - 1：VIDEO_LEVEL_1_0 - 2：VIDEO_LEVEL_1_1 - 3：VIDEO_LEVEL_1_2 - 4：VIDEO_LEVEL_1_3 - 5：VIDEO_LEVEL_2_0 - 6：VIDEO_LEVEL_2_1 - 7：VIDEO_LEVEL_2_2 - 8：VIDEO_LEVEL_3_0 - 9：VIDEO_LEVEL_3_1 - 10：VIDEO_LEVEL_3_2 - 11：VIDEO_LEVEL_4_0 - 12：VIDEO_LEVEL_4_1 - 13：VIDEO_LEVEL_4_2 - 14：VIDEO_LEVEL_5_0 - 15：VIDEO_LEVEL_5_1
	Level *int32 `json:"level,omitempty"`

	// 编码质量等级。  取值如下： - 1：VIDEO_PRESET_HSPEED2 - 2：VIDEO_PRESET_HSPEED - 3：VIDEO_PRESET_NORMAL > 值越大，表示编码的质量越高，转码耗时也越长。
	Preset *int32 `json:"preset,omitempty"`

	// 最大参考帧数。  取值范围： - H264：[1，8] - H265：固定值4  单位：帧。
	RefFramesCount *int32 `json:"ref_frames_count,omitempty"`

	// I帧最大间隔  取值范围：[2，10]。  默认值：5。  单位：秒。
	MaxIframesInterval *int32 `json:"max_iframes_interval,omitempty"`

	// 最大B帧间隔。  取值范围： - H264：[0，7]，默认值为4。 - H265：[0，7]，默认值为7。  单位：帧。
	BframesCount *int32 `json:"bframes_count,omitempty"`

	// 帧率。  取值范围：0或[5,60]之间的整数。  单位：帧每秒。  > 若设置的帧率不在取值范围内，则自动调整为0，若设置的帧率高于片源帧率，则自动调整为片源帧率。
	FrameRate *int32 `json:"frame_rate,omitempty"`

	// 视频宽度。  取值范围： - H.264：0或[32,4096]间2的倍数。 - H.265：0或[160,4096]间4的倍数。  单位：像素。  说明：若视频宽度设置为0，则视频宽度值自适应。
	Width *int32 `json:"width,omitempty"`

	// 视频高度。 - H.264：0或[32,2880]且必须为2的倍数。 - H.265：0或[96,2880]且必须为4的倍数。  单位：像素。  说明：若视频高度设置为0，则视频高度值自适应。
	Height *int32 `json:"height,omitempty"`

	// 黑边剪裁类型。  取值如下： - 0：不开启黑边剪裁。 - 1：开启黑边剪裁，低复杂度算法，针对长视频（>5分钟）。 - 2：开启黑边剪裁，高复杂度算法，针对短视频（<=5分钟）。
	BlackCut *int32 `json:"black_cut,omitempty"`
}

func (o Video) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Video struct{}"
	}

	return strings.Join([]string{"Video", string(data)}, " ")
}

type VideoOutputPolicy struct {
	value string
}

type VideoOutputPolicyEnum struct {
	TRANSCODE VideoOutputPolicy
	DISCARD   VideoOutputPolicy
	COPY      VideoOutputPolicy
}

func GetVideoOutputPolicyEnum() VideoOutputPolicyEnum {
	return VideoOutputPolicyEnum{
		TRANSCODE: VideoOutputPolicy{
			value: "transcode",
		},
		DISCARD: VideoOutputPolicy{
			value: "discard",
		},
		COPY: VideoOutputPolicy{
			value: "copy",
		},
	}
}

func (c VideoOutputPolicy) Value() string {
	return c.value
}

func (c VideoOutputPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *VideoOutputPolicy) UnmarshalJSON(b []byte) error {
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
