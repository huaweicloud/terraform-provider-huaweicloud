package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SecondarySourcesInfo 主入流信息
type SecondarySourcesInfo struct {

	// 频道源流URL，用于外部拉流
	Url *string `json:"url,omitempty"`

	// 码率。无需直播转码时，此参数为必填项  单位：bps。取值范围：0 - 104,857,600（100Mbps）
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 分辨率对应宽的值，非必填项  取值范围：0 - 4096（4K）
	Width *int32 `json:"width,omitempty"`

	// 分辨率对应高的值，非必填项  取值范围：0 - 2160（4K）
	Height *int32 `json:"height,omitempty"`

	// 是否使用bitrate来固定码率。默认值：false
	BitrateFor3u8 *bool `json:"bitrate_for3u8,omitempty"`

	// 协议为SRT_PUSH时的加密信息
	Passphrase *string `json:"passphrase,omitempty"`

	// 备入流地址列表
	BackupUrls *[]string `json:"backup_urls,omitempty"`

	// 协议为SRT_PULL时，拉流地址的streamid
	StreamId *string `json:"stream_id,omitempty"`

	// 协议为SRT_PULL时，拉流时延
	Latency *int32 `json:"latency,omitempty"`
}

func (o SecondarySourcesInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SecondarySourcesInfo struct{}"
	}

	return strings.Join([]string{"SecondarySourcesInfo", string(data)}, " ")
}
