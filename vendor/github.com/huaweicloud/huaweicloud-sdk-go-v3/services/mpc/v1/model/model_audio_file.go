package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AudioFile struct {

	// 音轨信息
	TracksInfo *[]TracksInfo `json:"tracks_info,omitempty"`

	Input *ObsObjInfo `json:"input,omitempty"`
}

func (o AudioFile) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AudioFile struct{}"
	}

	return strings.Join([]string{"AudioFile", string(data)}, " ")
}
