package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// This is a auto create Body Object
type StartTaskReq struct {

	// 源端节点AK（最大长度100个字符）。URL列表迁移任务不需要填写此参数。
	SrcAk *string `json:"src_ak,omitempty"`

	// 源端节点SK（最大长度100个字符）。URL列表迁移任务不需要填写此参数。
	SrcSk *string `json:"src_sk,omitempty"`

	// 源端节点临时Token
	SrcSecurityToken *string `json:"src_security_token,omitempty"`

	// 目的端节点AK（最大长度100个字符）。
	DstAk string `json:"dst_ak"`

	// 目的端节点SK（最大长度100个字符）。
	DstSk string `json:"dst_sk"`

	// 目标端节点临时Token
	DstSecurityToken *string `json:"dst_security_token,omitempty"`

	// CDN鉴权秘钥。
	SourceCdnAuthenticationKey *string `json:"source_cdn_authentication_key,omitempty"`

	// 迁移类型，标识是否为全量迁移，默认false（全量迁移）。 值为true时表示只重传失败对象。 值为空或者为false时表示全量迁移。
	MigrateFailedObject *bool `json:"migrate_failed_object,omitempty"`
}

func (o StartTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTaskReq struct{}"
	}

	return strings.Join([]string{"StartTaskReq", string(data)}, " ")
}
