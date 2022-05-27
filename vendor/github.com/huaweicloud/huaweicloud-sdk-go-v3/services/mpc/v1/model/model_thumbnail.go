package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Thumbnail struct {

	// 是否压缩抽帧图片生成tar包 - 0：表示压缩 - 1：表示不压缩
	Tar *int32 `json:"tar,omitempty"`

	Out *ObsObjInfo `json:"out,omitempty"`

	Params *ThumbnailPara `json:"params"`
}

func (o Thumbnail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Thumbnail struct{}"
	}

	return strings.Join([]string{"Thumbnail", string(data)}, " ")
}
