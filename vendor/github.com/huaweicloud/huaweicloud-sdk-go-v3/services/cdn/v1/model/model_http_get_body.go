package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HttpGetBody 证书配置查询响应体。
type HttpGetBody struct {

	// HTTPS证书是否启用，on：开启，off：关闭。
	HttpsStatus *string `json:"https_status,omitempty"`

	// 证书名字。
	CertificateName *string `json:"certificate_name,omitempty"`

	// HTTPS协议使用的证书内容，PEM编码格式。
	CertificateValue *string `json:"certificate_value,omitempty"`

	// 证书来源,1：华为云托管证书,0：自有证书。
	CertificateSource *int32 `json:"certificate_source,omitempty"`

	// 是否使用HTTP2.0，on：是，off：否。
	Http2Status *string `json:"http2_status,omitempty"`

	// 传输层安全性协议。
	TlsVersion *string `json:"tls_version,omitempty"`
}

func (o HttpGetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HttpGetBody struct{}"
	}

	return strings.Join([]string{"HttpGetBody", string(data)}, " ")
}
