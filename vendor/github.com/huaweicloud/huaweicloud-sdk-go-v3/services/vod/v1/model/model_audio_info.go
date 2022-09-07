package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AudioInfo struct {

	// 音频采样率(有效值范围)<br/> AUDIO_SAMPLE_AUTO (default), AUDIO_SAMPLE_22050：22050Hz<br/> AUDIO_SAMPLE_32000：32000Hz<br/> AUDIO_SAMPLE_44100：44100Hz<br/> AUDIO_SAMPLE_48000：48000Hz<br/> AUDIO_SAMPLE_96000：96000Hz<br/>
	SampleRate AudioInfoSampleRate `json:"sample_rate"`

	// 音频码率（单位：Kbps）<br/>
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 声道数(有效值范围)<br/> AUDIO_CHANNELS_1:单声道<br/> AUDIO_CHANNELS_2：双声道<br/> AUDIO_CHANNELS_5_1：5.1声道<br/>
	Channels AudioInfoChannels `json:"channels"`
}

func (o AudioInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AudioInfo struct{}"
	}

	return strings.Join([]string{"AudioInfo", string(data)}, " ")
}

type AudioInfoSampleRate struct {
	value string
}

type AudioInfoSampleRateEnum struct {
	AUDIO_SAMPLE_AUTO  AudioInfoSampleRate
	AUDIO_SAMPLE_22050 AudioInfoSampleRate
	AUDIO_SAMPLE_32000 AudioInfoSampleRate
	AUDIO_SAMPLE_44100 AudioInfoSampleRate
	AUDIO_SAMPLE_48000 AudioInfoSampleRate
	AUDIO_SAMPLE_96000 AudioInfoSampleRate
}

func GetAudioInfoSampleRateEnum() AudioInfoSampleRateEnum {
	return AudioInfoSampleRateEnum{
		AUDIO_SAMPLE_AUTO: AudioInfoSampleRate{
			value: "AUDIO_SAMPLE_AUTO",
		},
		AUDIO_SAMPLE_22050: AudioInfoSampleRate{
			value: "AUDIO_SAMPLE_22050",
		},
		AUDIO_SAMPLE_32000: AudioInfoSampleRate{
			value: "AUDIO_SAMPLE_32000",
		},
		AUDIO_SAMPLE_44100: AudioInfoSampleRate{
			value: "AUDIO_SAMPLE_44100",
		},
		AUDIO_SAMPLE_48000: AudioInfoSampleRate{
			value: "AUDIO_SAMPLE_48000",
		},
		AUDIO_SAMPLE_96000: AudioInfoSampleRate{
			value: "AUDIO_SAMPLE_96000",
		},
	}
}

func (c AudioInfoSampleRate) Value() string {
	return c.value
}

func (c AudioInfoSampleRate) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AudioInfoSampleRate) UnmarshalJSON(b []byte) error {
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

type AudioInfoChannels struct {
	value string
}

type AudioInfoChannelsEnum struct {
	AUDIO_CHANNELS_1   AudioInfoChannels
	AUDIO_CHANNELS_2   AudioInfoChannels
	AUDIO_CHANNELS_5_1 AudioInfoChannels
}

func GetAudioInfoChannelsEnum() AudioInfoChannelsEnum {
	return AudioInfoChannelsEnum{
		AUDIO_CHANNELS_1: AudioInfoChannels{
			value: "AUDIO_CHANNELS_1",
		},
		AUDIO_CHANNELS_2: AudioInfoChannels{
			value: "AUDIO_CHANNELS_2",
		},
		AUDIO_CHANNELS_5_1: AudioInfoChannels{
			value: "AUDIO_CHANNELS_5_1",
		},
	}
}

func (c AudioInfoChannels) Value() string {
	return c.value
}

func (c AudioInfoChannels) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AudioInfoChannels) UnmarshalJSON(b []byte) error {
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
