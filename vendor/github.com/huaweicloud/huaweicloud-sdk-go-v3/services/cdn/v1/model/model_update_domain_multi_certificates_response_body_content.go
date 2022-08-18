package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateDomainMultiCertificatesResponseBodyContent struct {

	// 域名列表
	DomainName string `json:"domain_name"`

	// https开关(0：\"关闭\"；1：\"设置证书\")
	HttpsSwitch *int32 `json:"https_switch,omitempty"`

	// 回源方式:1：\"PROTOCOL_FOLLOW\"；2：\"HTTP\"(默认)，3：\"https\"
	AccessOriginWay *int32 `json:"access_origin_way,omitempty"`

	// 强制跳转HTTPS（0：不强制；1：强制） 为空值时默认设置为关闭。（建议使用force_redirect_config修改配置）
	ForceRedirectHttps *int32 `json:"force_redirect_https,omitempty"`

	ForceRedirectConfig *ForceRedirect `json:"force_redirect_config,omitempty"`

	// http2.0（0：关闭；1：开启）
	Http2 *int32 `json:"http2,omitempty"`

	// 证书名称。（长度限制为3-32字符）。
	CertName *string `json:"cert_name,omitempty"`

	// 证书内容
	Certificate *string `json:"certificate,omitempty"`

	// 证书类型（0为自有证书 ， 1为托管证书）
	CertificateType *int32 `json:"certificate_type,omitempty"`

	// 证书过期时间
	ExpirationTime *int64 `json:"expiration_time,omitempty"`
}

func (o UpdateDomainMultiCertificatesResponseBodyContent) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainMultiCertificatesResponseBodyContent struct{}"
	}

	return strings.Join([]string{"UpdateDomainMultiCertificatesResponseBodyContent", string(data)}, " ")
}
