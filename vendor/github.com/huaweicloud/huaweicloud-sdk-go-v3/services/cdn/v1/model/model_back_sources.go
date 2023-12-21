package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BackSources 高级回源信息源站配置。
type BackSources struct {

	// 源站类型, ipaddr：源站IP，domain：源站域名，obs_bucket：OBS桶域名。
	SourcesType string `json:"sources_type"`

	// 源站IP或者域名。
	IpOrDomain string `json:"ip_or_domain"`

	// OBS桶类型： - “private”， 私有桶： - “public”，公有桶。
	ObsBucketType *string `json:"obs_bucket_type,omitempty"`

	// HTTP端口，取值范围：1-65535。
	HttpPort *int32 `json:"http_port,omitempty"`

	// HTTPS端口，取值范围：1-65535。
	HttpsPort *int32 `json:"https_port,omitempty"`
}

func (o BackSources) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackSources struct{}"
	}

	return strings.Join([]string{"BackSources", string(data)}, " ")
}
