package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HttpPutBody 设置证书请求体
type HttpPutBody struct {

	// HTTPS证书是否启用，on：开启，off：关闭。
	HttpsStatus *string `json:"https_status,omitempty"`

	// 证书名字，长度限制为3-64字符。  > 当证书开启时必传。
	CertificateName *string `json:"certificate_name,omitempty"`

	// HTTPS协议使用的证书内容，当证书开启时必传。  > PEM编码格式。
	CertificateValue *string `json:"certificate_value,omitempty"`

	// HTTPS协议使用的私钥，当证书开启时必传。  > PEM编码格式。
	PrivateKey *string `json:"private_key,omitempty"`

	// 证书来源,0：自有证书。  > 证书开启时必传
	CertificateSource *int32 `json:"certificate_source,omitempty"`

	// 是否使用HTTP2.0，on：是，off：否。  > 默认关闭，https_status=off时，该值不生效。
	Http2Status *string `json:"http2_status,omitempty"`

	// 传输层安全性协议， 目前支持TLSv1.0/1.1/1.2/1.3四个版本的协议，CDN默认开启TLS1.1/1.2/1.3，不可全部关闭。  > 1.需开启连续或单个版本号，例：不可仅开启TLS1.0/1.2而关闭TLS1.1。  > 2.多版本开启时，使用逗号拼接传输，例：TLSv1.1,TLSv1.2。
	TlsVersion *string `json:"tls_version,omitempty"`
}

func (o HttpPutBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HttpPutBody struct{}"
	}

	return strings.Join([]string{"HttpPutBody", string(data)}, " ")
}
