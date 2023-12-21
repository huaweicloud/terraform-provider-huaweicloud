package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InheritConfigQuery 鉴权继承，为M3U8/MPD索引文件下的TS/MP4文件添加鉴权参数，解决因鉴权不通过导致的TS/MP4文件无法播放的问题。
type InheritConfigQuery struct {

	// 是否开启鉴权继承，on：开启,off：关闭。
	Status string `json:"status"`

	// 鉴权继承配置， m3u8：M3U8,mpd：MPD,“m3u8,mpd”。
	InheritType *string `json:"inherit_type,omitempty"`

	// 鉴权继承的文件类型时间, sys_time：当前系统时间，parent_url_time：与m3u8和mpd访问链接保持一致。
	InheritTimeType *string `json:"inherit_time_type,omitempty"`
}

func (o InheritConfigQuery) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InheritConfigQuery struct{}"
	}

	return strings.Join([]string{"InheritConfigQuery", string(data)}, " ")
}
