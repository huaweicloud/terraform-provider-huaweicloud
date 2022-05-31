package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateMergeChannelsReq struct {
	MultiAudio *MpcMultiAudio `json:"multi_audio,omitempty"`
}

func (o CreateMergeChannelsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMergeChannelsReq struct{}"
	}

	return strings.Join([]string{"CreateMergeChannelsReq", string(data)}, " ")
}
