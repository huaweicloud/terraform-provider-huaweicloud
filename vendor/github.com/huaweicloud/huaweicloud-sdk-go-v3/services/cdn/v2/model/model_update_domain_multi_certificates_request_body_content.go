package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainMultiCertificatesRequestBodyContent https配置。
type UpdateDomainMultiCertificatesRequestBodyContent struct {

	// 域名列表,逗号分割，上限50个域名
	DomainName string `json:"domain_name"`

	// https开关（0：\"关闭\"；1：\"设置证书\" https_switch为1时，证书参数不能为空）
	HttpsSwitch int32 `json:"https_switch"`

	// 回源方式:1：\"回源跟随\"；2：\"http\"(默认)，3：\"https\"  为空值时默认设置为http
	AccessOriginWay *int32 `json:"access_origin_way,omitempty"`

	// 强制跳转HTTPS（0：不强制；1：强制） 为空值时默认设置为关闭。（此参数即将下线，建议使用force_redirect_config修改配置）
	ForceRedirectHttps *int32 `json:"force_redirect_https,omitempty"`

	ForceRedirectConfig *ForceRedirect `json:"force_redirect_config,omitempty"`

	// http2.0（0：关闭；1：开启） 为空值时默认设置为关闭
	Http2 *int32 `json:"http2,omitempty"`

	// 证书名称（设置证书必填）（长度限制为3-64字符）。
	CertName *string `json:"cert_name,omitempty"`

	// HTTPS协议使用的SSL证书内容，仅支持PEM编码格式。不启用证书则无需输入。初次配置证书时必传。
	Certificate *string `json:"certificate,omitempty"`

	// HTTPS协议使用的SSL证书私钥内容，仅支持PEM编码格式。不启用证书则无需输入。初次配置证书时必传。
	PrivateKey *string `json:"private_key,omitempty"`

	// 证书类型（0为自有证书；2为SCM证书；不传默认为自有证书）
	CertificateType *int32 `json:"certificate_type,omitempty"`

	// SCM证书id
	ScmCertificateId *string `json:"scm_certificate_id,omitempty"`
}

func (o UpdateDomainMultiCertificatesRequestBodyContent) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainMultiCertificatesRequestBodyContent struct{}"
	}

	return strings.Join([]string{"UpdateDomainMultiCertificatesRequestBodyContent", string(data)}, " ")
}
