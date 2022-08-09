package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 源站配置。
type SourcesConfig struct {

	// 源站IP或者域名。
	OriginAddr string `json:"origin_addr"`

	// 源站类型（\"ipaddr\"：\"IP源站\"，\"domain\"：\"域名源站\"，\"obs_bucket\"：\"OBS Bucket源站\"）。
	OriginType string `json:"origin_type"`

	// 源站优先级（70：主，30：备）。
	Priority int32 `json:"priority"`

	// 是否开启Obs静态网站托管，源站类型为obs_bucket时传递(off：关闭，on：开启)。
	ObsWebHostingStatus *string `json:"obs_web_hosting_status,omitempty"`

	// HTTP端口，默认80。
	HttpPort *int32 `json:"http_port,omitempty"`

	// HTTPS端口，默认443。
	HttpsPort *int32 `json:"https_port,omitempty"`

	// 回源HOST，默认加速域名。
	HostName *string `json:"host_name,omitempty"`
}

func (o SourcesConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourcesConfig struct{}"
	}

	return strings.Join([]string{"SourcesConfig", string(data)}, " ")
}
