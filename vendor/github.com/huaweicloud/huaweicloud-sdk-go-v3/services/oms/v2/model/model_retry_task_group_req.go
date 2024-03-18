package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryTaskGroupReq 重试迁移任务组请求参数
type RetryTaskGroupReq struct {

	// 源端ak（最大长度100个字符）
	SrcAk *string `json:"src_ak,omitempty"`

	// 源端sk（最大长度100个字符）
	SrcSk *string `json:"src_sk,omitempty"`

	// 用于谷歌云Cloud Storage鉴权
	JsonAuthFile *string `json:"json_auth_file,omitempty"`

	// 目的端ak（最大长度100个字符）
	DstAk *string `json:"dst_ak,omitempty"`

	// 目的端sk（最大长度100个字符）
	DstSk *string `json:"dst_sk,omitempty"`

	// cdn鉴权秘钥
	SourceCdnAuthenticationKey *string `json:"source_cdn_authentication_key,omitempty"`

	// 失败任务重试方式，标识是否为全量重新迁移，默认false（全量重新迁移）。 值为true时表示只重传失败对象。 值为空或者为false时表示全量重新迁移（默认跳过目的端已迁移对象）。
	MigrateFailedObject *bool `json:"migrate_failed_object,omitempty"`
}

func (o RetryTaskGroupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryTaskGroupReq struct{}"
	}

	return strings.Join([]string{"RetryTaskGroupReq", string(data)}, " ")
}
