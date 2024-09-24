package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CertificatesGetBody 双证书配置查询响应体。
type CertificatesGetBody struct {

	// 证书来源,0：自有证书。2：SCM证书。
	CertificateSource *int32 `json:"certificate_source,omitempty"`

	// SCM证书id
	ScmCertificateId *string `json:"scm_certificate_id,omitempty"`

	// 证书类型，server：国际证书；server_sm：国密证书。
	CertificateType *string `json:"certificate_type,omitempty"`

	// 证书名字。
	CertificateName *string `json:"certificate_name,omitempty"`

	// HTTPS协议使用的证书内容，PEM编码格式。
	CertificateValue *string `json:"certificate_value,omitempty"`

	// 国密证书加密证书内容，PEM编码格式。
	EncCertificateValue *string `json:"enc_certificate_value,omitempty"`

	// 证书过期时间。  > UTC时间。
	ExpireTime *int64 `json:"expire_time,omitempty"`
}

func (o CertificatesGetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CertificatesGetBody struct{}"
	}

	return strings.Join([]string{"CertificatesGetBody", string(data)}, " ")
}
