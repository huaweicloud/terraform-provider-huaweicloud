package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VideoSeek 视频拖拽配置。 > 1. 需同步开启FLV、MP4格式文件的URL参数功能，并选择“忽略参数”。 > 2. 关闭视频拖拽功能时，FLV时间拖拽功能失效。
type VideoSeek struct {

	// 视频拖拽开关， true：开启，false：关闭  > 当本字段设置为“false”时，查询域名配置接口将不会返回视频拖拽配置信息。
	EnableVideoSeek bool `json:"enable_video_seek"`

	// flv时间拖拽开关， true：开启，false：关闭。
	EnableFlvByTimeSeek *bool `json:"enable_flv_by_time_seek,omitempty"`

	// 自定义用户请求URL中视频播放的开始参数，支持使用数字0-9、字符a-z、A-Z，及\"_\"，长度≤64个字符。
	StartParameter *string `json:"start_parameter,omitempty"`

	// 自定义用户请求URL中视频播放的结束参数，支持使用数字0-9、字符a-z、A-Z，及\"_\"，长度≤64个字符。
	EndParameter *string `json:"end_parameter,omitempty"`
}

func (o VideoSeek) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoSeek struct{}"
	}

	return strings.Join([]string{"VideoSeek", string(data)}, " ")
}
