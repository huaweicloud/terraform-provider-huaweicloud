package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ObsFileAddr struct {

	// OBS的bucket名称
	Bucket string `json:"bucket"`

	// OBS Bucket所在的区域，且必须保持与使用的直播服务区域保持一致。
	Location string `json:"location"`

	// OBS对象路径，遵守OSS Object定义 - 当用于指示input时，需要指定到具体对象 - 当用于指示output时，只需指定到转码结果期望存放的路径
	Object string `json:"object"`
}

func (o ObsFileAddr) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsFileAddr struct{}"
	}

	return strings.Join([]string{"ObsFileAddr", string(data)}, " ")
}
