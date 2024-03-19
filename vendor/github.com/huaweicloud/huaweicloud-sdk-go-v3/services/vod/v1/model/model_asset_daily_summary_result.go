package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AssetDailySummaryResult struct {

	// 播放日期，格式为：yyyyMMdd000000。
	Date *string `json:"date,omitempty"`

	// 日播放统计数据文件的下载地址，有效期为12个小时。  文件内容格式：[域名]\\t[媒资ID]\\t[日期]\\t[播放流量]\\t[播放次数]  播放次数统计说明： - HLS文件：访问M3U8文件时会统计播放次数，访问TS文件时不会统计播放次数。 - 其它文件：如MP4文件，当播放请求带有range且range的start参数不等于0时，不统计播放次数。其它情况下，会统计播放次数。
	Link *string `json:"link,omitempty"`
}

func (o AssetDailySummaryResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssetDailySummaryResult struct{}"
	}

	return strings.Join([]string{"AssetDailySummaryResult", string(data)}, " ")
}
