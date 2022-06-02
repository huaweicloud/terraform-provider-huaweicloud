package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Audio struct {

	// 输出策略。  取值如下： - discard - transcode  >- 当视频参数中的“output_policy”为\"discard\"，且音频参数中的“output_policy”为“transcode”时，表示只输出音频。 >- 当视频参数中的“output_policy”为\"transcode\"，且音频参数中的“output_policy”为“discard”时，表示只输出视频。 >- 同时为\"discard\"时不合法。 >- 同时为“transcode”时，表示输出音视频。
	OutputPolicy *AudioOutputPolicy `json:"output_policy,omitempty"`

	// 音频编码格式。  取值如下：  - 1：AAC格式。 - 2：HEAAC1格式 。 - 3：HEAAC2格式。 - 4：MP3格式 。
	Codec int32 `json:"codec"`

	// 音频采样率。  取值如下：  - 1：AUDIO_SAMPLE_AUTO - 2：AUDIO_SAMPLE_22050（22050Hz） - 3：AUDIO_SAMPLE_32000（32000Hz） - 4：AUDIO_SAMPLE_44100（44100Hz） - 5：AUDIO_SAMPLE_48000（48000Hz） - 6：AUDIO_SAMPLE_96000（96000Hz）
	SampleRate int32 `json:"sample_rate"`

	// 音频码率。  取值范围：0或[8,1000]。  单位：kbit/s。
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 声道数。  取值如下： - 1：AUDIO_CHANNELS_1 - 2：AUDIO_CHANNELS_2 - 6：AUDIO_CHANNELS_5_1
	Channels int32 `json:"channels"`
}

func (o Audio) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Audio struct{}"
	}

	return strings.Join([]string{"Audio", string(data)}, " ")
}

type AudioOutputPolicy struct {
	value string
}

type AudioOutputPolicyEnum struct {
	TRANSCODE AudioOutputPolicy
	DISCARD   AudioOutputPolicy
	COPY      AudioOutputPolicy
}

func GetAudioOutputPolicyEnum() AudioOutputPolicyEnum {
	return AudioOutputPolicyEnum{
		TRANSCODE: AudioOutputPolicy{
			value: "transcode",
		},
		DISCARD: AudioOutputPolicy{
			value: "discard",
		},
		COPY: AudioOutputPolicy{
			value: "copy",
		},
	}
}

func (c AudioOutputPolicy) Value() string {
	return c.value
}

func (c AudioOutputPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AudioOutputPolicy) UnmarshalJSON(b []byte) error {
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
