package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PicInfo struct {

	// 截图文件名。
	PicName *string `json:"pic_name,omitempty"`
}

func (o PicInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PicInfo struct{}"
	}

	return strings.Join([]string{"PicInfo", string(data)}, " ")
}
