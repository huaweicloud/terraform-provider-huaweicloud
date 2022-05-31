package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 截图结果
type ThumbnailRsp struct {

	// 截图在视频中的时间偏移，单位为秒。
	Offset int32 `json:"offset"`

	// 截图访问URL
	Url string `json:"url"`
}

func (o ThumbnailRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ThumbnailRsp struct{}"
	}

	return strings.Join([]string{"ThumbnailRsp", string(data)}, " ")
}
