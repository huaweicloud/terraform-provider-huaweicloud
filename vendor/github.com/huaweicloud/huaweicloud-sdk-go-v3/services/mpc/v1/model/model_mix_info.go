package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MixInfo struct {

	// 合成任务原始视频配置
	Inputs *[]InputSetting `json:"inputs,omitempty"`

	Layout *MixInfoLayout `json:"layout,omitempty"`
}

func (o MixInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MixInfo struct{}"
	}

	return strings.Join([]string{"MixInfo", string(data)}, " ")
}
