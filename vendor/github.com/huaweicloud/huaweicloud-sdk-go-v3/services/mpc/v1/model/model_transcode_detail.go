package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TranscodeDetail struct {

	// 一进多出情况下部分转码失败的情况。
	MultitaskInfo *[]MultiTaskInfo `json:"multitask_info,omitempty"`

	InputFile *SourceInfo `json:"input_file,omitempty"`
}

func (o TranscodeDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TranscodeDetail struct{}"
	}

	return strings.Join([]string{"TranscodeDetail", string(data)}, " ")
}
