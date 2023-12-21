package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartSyncTaskReq 启动同步任务body体
type StartSyncTaskReq struct {

	// 源端节点AK（最大长度100个字符）。URL列表迁移任务不需要填写此参数。
	SrcAk string `json:"src_ak"`

	// 源端节点SK（最大长度100个字符）。URL列表迁移任务不需要填写此参数。
	SrcSk string `json:"src_sk"`

	// 目的端节点AK（最大长度100个字符）。
	DstAk string `json:"dst_ak"`

	// 目的端节点SK（最大长度100个字符）。
	DstSk string `json:"dst_sk"`

	// CDN鉴权秘钥。
	SourceCdnAuthenticationKey *string `json:"source_cdn_authentication_key,omitempty"`
}

func (o StartSyncTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartSyncTaskReq struct{}"
	}

	return strings.Join([]string{"StartSyncTaskReq", string(data)}, " ")
}
