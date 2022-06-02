package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type InputSetting struct {
	Input *ObsObjInfo `json:"input"`

	// 原视频的id,为整数类型数值字符串。用于匹配后面的布局配置。
	PaneId string `json:"pane_id"`

	// 该视频采取的音频策略。DISCARD表示丢弃音频，合成的视频中不会出现该视频的音频。 RESERVE表示保留音频，合成的视频中会出现该视频音频。如果多个原视频配置了RESERVE，合成的视频文件的音频，是多个原 视频音频的混合。默认会丢弃音频。
	AudioPolicy *InputSettingAudioPolicy `json:"audio_policy,omitempty"`
}

func (o InputSetting) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InputSetting struct{}"
	}

	return strings.Join([]string{"InputSetting", string(data)}, " ")
}

type InputSettingAudioPolicy struct {
	value string
}

type InputSettingAudioPolicyEnum struct {
	DISCARD InputSettingAudioPolicy
	RESERVE InputSettingAudioPolicy
}

func GetInputSettingAudioPolicyEnum() InputSettingAudioPolicyEnum {
	return InputSettingAudioPolicyEnum{
		DISCARD: InputSettingAudioPolicy{
			value: "DISCARD",
		},
		RESERVE: InputSettingAudioPolicy{
			value: "RESERVE",
		},
	}
}

func (c InputSettingAudioPolicy) Value() string {
	return c.value
}

func (c InputSettingAudioPolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *InputSettingAudioPolicy) UnmarshalJSON(b []byte) error {
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
