package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InheritConfig 鉴权继承，为M3U8/MPD索引文件下的TS/MP4文件添加鉴权参数，解决因鉴权不通过导致的TS/MP4文件无法播放的问题。
type InheritConfig struct {

	// 是否开启鉴权继承，on：开启,off：关闭。
	Status string `json:"status"`

	// 鉴权继承配置， m3u8：M3U8,mpd：MPD,输入多个时用“,”分开，例如“m3u8,mpd”  > 开启鉴权继承时，该参数必填。
	InheritType *string `json:"inherit_type,omitempty"`

	// 鉴权继承的文件类型时间, sys_time：当前系统时间，parent_url_time：与m3u8和mpd访问链接保持一致。  > 开启鉴权继承时，该参数必填。
	InheritTimeType *string `json:"inherit_time_type,omitempty"`
}

func (o InheritConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InheritConfig struct{}"
	}

	return strings.Join([]string{"InheritConfig", string(data)}, " ")
}
