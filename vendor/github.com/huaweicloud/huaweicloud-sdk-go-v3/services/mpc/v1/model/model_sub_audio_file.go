package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SubAudioFile struct {

	// 音轨信息
	TracksInfo *[]TracksInfo `json:"tracks_info,omitempty"`

	Input *ObsObjInfo `json:"input,omitempty"`

	Output *ObsObjInfo `json:"output,omitempty"`

	// 输出文件名。
	OutputFilename *string `json:"output_filename,omitempty"`
}

func (o SubAudioFile) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SubAudioFile struct{}"
	}

	return strings.Join([]string{"SubAudioFile", string(data)}, " ")
}
