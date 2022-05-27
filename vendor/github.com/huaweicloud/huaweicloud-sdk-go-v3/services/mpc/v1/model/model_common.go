package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Common struct {

	// 是否开启高清低码功能。  取值如下： - false：关闭。 - true：开启。
	Pvc bool `json:"PVC"`

	// HLS分片间隔，仅封装类型“pack_type”取值为1或3时，该参数生效。  取值范围：[2，10]。  单位：秒。
	HlsInterval int32 `json:"hls_interval"`

	// DASH间隔，仅封装类型“pack_type”取值为2或3时，该参数生效。  取值范围：[2，10]。  单位：秒。
	DashInterval int32 `json:"dash_interval"`

	// 封装类型。  取值如下： - 1：HLS - 2：DASH - 3：HLS+DASH - 4：MP4 - 5：MP3 - 6：ADTS  > pack_type设置为5和6时，不能设置视频参数。
	PackType int32 `json:"pack_type"`
}

func (o Common) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Common struct{}"
	}

	return strings.Join([]string{"Common", string(data)}, " ")
}
