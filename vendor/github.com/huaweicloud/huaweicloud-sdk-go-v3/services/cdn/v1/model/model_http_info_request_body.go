package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HttpInfoRequestBody struct {

	// 证书名字。（长度限制为3-32字符）。
	CertName string `json:"cert_name"`

	// HTTPS证书是否启用。0：不启用，此时无需填写证书及私钥参数；1：启用HTTPS加速并协议跟随回源；2：启用HTTPS加速并HTTP回源；3：启用HTTPS加速并HTTPS回源，首次配置证书需要传递证书及私钥，如已有证书可不用传证书及私钥。
	HttpsStatus int32 `json:"https_status"`

	// 功能说明：HTTPS协议使用的证书内容，不启用证书则无需输入。取值范围：PEM编码格式。初次配置证书时必传。
	Certificate *string `json:"certificate,omitempty"`

	// 功能说明： HTTPS协议使用的私钥，不启用证书则无需输入。取值范围：PEM编码格式。初次配置证书时必传。
	PrivateKey *string `json:"private_key,omitempty"`

	// 是否使用HTTP2.0。（1：是，0：否。）
	Http2 *int32 `json:"http2,omitempty"`

	// 证书类型。1：代表华为云托管证书；0：表示自有证书。 默认值0。
	CertificateType *int32 `json:"certificate_type,omitempty"`

	// 强制跳转HTTPS（0：不强制；1：强制） 为空值时默认设置为关闭。（建议使用force_redirect_config修改配置）
	ForceRedirectHttps *int32 `json:"force_redirect_https,omitempty"`

	ForceRedirectConfig *ForceRedirect `json:"force_redirect_config,omitempty"`
}

func (o HttpInfoRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HttpInfoRequestBody struct{}"
	}

	return strings.Join([]string{"HttpInfoRequestBody", string(data)}, " ")
}
