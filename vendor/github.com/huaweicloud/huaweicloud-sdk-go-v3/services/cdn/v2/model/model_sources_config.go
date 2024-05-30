package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SourcesConfig 源站配置。
type SourcesConfig struct {

	// 源站类型， - ipaddr：源站IP； - domain：源站域名； - obs_bucket：OBS桶域名； - third_bucket：第三方桶。
	OriginType string `json:"origin_type"`

	// 源站IP或者域名。
	OriginAddr string `json:"origin_addr"`

	// 源站优先级，70：主，30：备。
	Priority int32 `json:"priority"`

	// 权重，取值范围1-100。
	Weight *int32 `json:"weight,omitempty"`

	// 是否开启OBS静态网站托管，源站类型为obs_bucket时传递，off：关闭，on：开启。
	ObsWebHostingStatus *string `json:"obs_web_hosting_status,omitempty"`

	// HTTP端口，默认80,端口取值取值范围1-65535。
	HttpPort *int32 `json:"http_port,omitempty"`

	// HTTPS端口，默认443,端口取值取值范围1-65535。
	HttpsPort *int32 `json:"https_port,omitempty"`

	// 回源HOST，默认加速域名。
	HostName *string `json:"host_name,omitempty"`

	// OBS桶类型，源站类型是“OBS桶域名”时需要传该参数，不传默认为“public”。   - private: 私有桶（除桶ACL授权外的其他用户无桶的访问权限）。   - public: 公有桶（任何用户都可以对桶内对象进行读操作）。
	ObsBucketType *string `json:"obs_bucket_type,omitempty"`

	// 第三方对象存储访问密钥。  > 源站类型为第三方桶时必填
	BucketAccessKey *string `json:"bucket_access_key,omitempty"`

	// 第三方对象存储密钥。  > 源站类型为第三方桶时必填
	BucketSecretKey *string `json:"bucket_secret_key,omitempty"`

	// 第三方对象存储区域。  > 源站类型为第三方桶时必填
	BucketRegion *string `json:"bucket_region,omitempty"`

	// 第三方对象存储名称。  > 源站类型为第三方桶时必填
	BucketName *string `json:"bucket_name,omitempty"`
}

func (o SourcesConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourcesConfig struct{}"
	}

	return strings.Join([]string{"SourcesConfig", string(data)}, " ")
}
