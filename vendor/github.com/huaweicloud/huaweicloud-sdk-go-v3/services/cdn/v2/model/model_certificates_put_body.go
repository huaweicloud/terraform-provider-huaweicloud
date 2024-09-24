package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CertificatesPutBody 配置双证书时必传，需要同时传入国际证书和国密证书，不支持传两个国际证书或两个国密证书。   > - 您也可以在certificates参数下传入一个国际证书或一个国密证书。   > - 如果certificates传了证书（国际证书、国密证书或国际+国密双证书），外层证书配置将失效，仅保留当前参数传入的证书信息。
type CertificatesPutBody struct {

	// 证书来源，0：自有证书。2：SCM证书。
	CertificateSource *int32 `json:"certificate_source,omitempty"`

	// SCM证书id
	ScmCertificateId *string `json:"scm_certificate_id,omitempty"`

	// 证书类型，server：国际证书；server_sm：国密证书。
	CertificateType string `json:"certificate_type"`

	// 证书名字，长度限制为3-64字符。
	CertificateName string `json:"certificate_name"`

	// HTTPS协议使用的证书内容。  > PEM编码格式。
	CertificateValue string `json:"certificate_value"`

	// HTTPS协议使用的私钥。  > PEM编码格式。
	PrivateKey string `json:"private_key"`

	// 加密证书内容，证书类型为国密证书时必传。  > PEM编码格式。
	EncCertificateValue *string `json:"enc_certificate_value,omitempty"`

	// 加密私钥内容，证书类型为国密证书时必传。  > PEM编码格式。
	EncPrivateKey *string `json:"enc_private_key,omitempty"`
}

func (o CertificatesPutBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CertificatesPutBody struct{}"
	}

	return strings.Join([]string{"CertificatesPutBody", string(data)}, " ")
}
