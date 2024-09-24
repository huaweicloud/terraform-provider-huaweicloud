package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GmCertificateInfo struct {

	// 证书来源， 可选,  scm: 云证书管理服务，user：默认，用户自有
	Source *string `json:"source,omitempty"`

	// SCM证书名， 可选
	CertName *string `json:"cert_name,omitempty"`

	// SCM证书ID, 证书来源为scm时必选
	CertId *string `json:"cert_id,omitempty"`

	// 国密签名证书内容
	SignCertificate *string `json:"sign_certificate,omitempty"`

	// 国密签名私钥内容
	SignCertificateKey *string `json:"sign_certificate_key,omitempty"`

	// 国密加密证书内容
	EncCertificate *string `json:"enc_certificate,omitempty"`

	// 国密加密私钥内容
	EncCertificateKey *string `json:"enc_certificate_key,omitempty"`
}

func (o GmCertificateInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GmCertificateInfo struct{}"
	}

	return strings.Join([]string{"GmCertificateInfo", string(data)}, " ")
}
