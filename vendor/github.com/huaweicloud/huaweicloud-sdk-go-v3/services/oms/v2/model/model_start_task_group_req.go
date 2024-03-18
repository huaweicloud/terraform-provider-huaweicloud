package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartTaskGroupReq This is a auto create Body Object
type StartTaskGroupReq struct {

	// 源端节点AK（最大长度100个字符）。URL列表迁移任务不需要填写此参数。
	SrcAk *string `json:"src_ak,omitempty"`

	// 源端节点SK（最大长度100个字符）。URL列表迁移任务不需要填写此参数。
	SrcSk *string `json:"src_sk,omitempty"`

	// 用于谷歌云Cloud Storage鉴权
	JsonAuthFile *string `json:"json_auth_file,omitempty"`

	// 目的端节点AK（最大长度100个字符）。
	DstAk string `json:"dst_ak"`

	// 目的端节点SK（最大长度100个字符）。
	DstSk string `json:"dst_sk"`

	// CDN鉴权秘钥。
	SourceCdnAuthenticationKey *string `json:"source_cdn_authentication_key,omitempty"`
}

func (o StartTaskGroupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTaskGroupReq struct{}"
	}

	return strings.Join([]string{"StartTaskGroupReq", string(data)}, " ")
}
