package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Crop struct {

	// 截取的视频时长。  单位：秒  从0秒开始算起
	Duration *int32 `json:"duration,omitempty"`
}

func (o Crop) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Crop struct{}"
	}

	return strings.Join([]string{"Crop", string(data)}, " ")
}
