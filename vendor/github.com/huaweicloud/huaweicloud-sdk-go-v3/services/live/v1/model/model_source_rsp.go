package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SourceRsp RTMP推流地址。只有频道入流协议为RTMP_PUSH时，会返回RTMP推流地址
type SourceRsp struct {

	// RTMP推流地址
	Url *string `json:"url,omitempty"`

	// 码率。  单位：bps。取值范围：0 - 104,857,600（100Mbps）
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 分辨率对应宽的值。取值范围：0 - 4096（4K）
	Width *int32 `json:"width,omitempty"`

	// 分辨率对应高的值。取值范围：0 - 2160（4K）
	Height *int32 `json:"height,omitempty"`

	// 描述是否使用该流做截图
	EnableSnapshot *bool `json:"enable_snapshot,omitempty"`
}

func (o SourceRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceRsp struct{}"
	}

	return strings.Join([]string{"SourceRsp", string(data)}, " ")
}
