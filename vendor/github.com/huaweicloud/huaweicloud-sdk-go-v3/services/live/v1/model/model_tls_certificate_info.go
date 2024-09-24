package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TlsCertificateInfo struct {

	// 证书来源， 可选,  scm: 云证书管理服务，user：默认，用户自有
	Source *string `json:"source,omitempty"`

	// SCM证书名， 证书来源为scm时可选
	CertName *string `json:"cert_name,omitempty"`

	// SCM证书ID, 证书来源为scm时必选
	CertId *string `json:"cert_id,omitempty"`

	// 证书内容，证书来源为user时必选
	Certificate *string `json:"certificate,omitempty"`

	// 私钥内容，证书来源为user时必选
	CertificateKey *string `json:"certificate_key,omitempty"`
}

func (o TlsCertificateInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TlsCertificateInfo struct{}"
	}

	return strings.Join([]string{"TlsCertificateInfo", string(data)}, " ")
}
